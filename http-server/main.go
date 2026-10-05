package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Event struct {
	UserID   string    `json:"user_id"`
	Action   string    `json:"action"`
	Value    int64     `json:"value"`
	Received time.Time `json:"received"`
}


type UserStat struct {
	UserID string `json:"user_id"`
	Count  int64  `json:"count"`
	Value  int64  `json:"value"`
}

type shard struct {
	mu   sync.RWMutex
	data map[string]*UserStat
}

type snapshot struct {
	Processed int64      `json:"processed"`
	Dropped   int64      `json:"dropped"`
	Queued    int        `json:"queued"`
	Workers   int        `json:"workers"`
	Uptime    string     `json:"uptime"`
	TopUsers  []UserStat `json:"top_users"`
}

type broker struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

func newBroker() *broker {
	return &broker{
		clients: make(map[chan []byte]struct{}),
	}
}

func (b *broker) subscribe() chan []byte {
	ch := make(chan []byte, 16)

	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()

	return ch
}

func (b *broker) unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *broker) publish(msg []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- msg:
		default:
		}
	}
}

type server struct {
	startedAt time.Time
	queue     chan Event
	shards    []shard
	streams   *broker

	processed atomic.Int64
	dropped   atomic.Int64
}

func newServer(workerCount int, queueSize int, shardCount int) *server {
	s := &server{
		startedAt: time.Now(),
		queue:     make(chan Event, queueSize),
		shards:    make([]shard, shardCount),
		streams:   newBroker(),
	}

	for i := range s.shards {
		s.shards[i].data = make(map[string]*UserStat)
	}

	for i := 0; i < workerCount; i++ {
		go s.worker()
	}

	return s
}

func (s *server) worker() {
	for event := range s.queue {
		s.apply(event)
		s.processed.Add(1)
	}
}

func (s *server) apply(event Event) {
	idx := shardIndex(event.UserID, len(s.shards))
	sh := &s.shards[idx]

	sh.mu.Lock()
	stat := sh.data[event.UserID]
	if stat == nil {
		stat = &UserStat{UserID: event.UserID}
		sh.data[event.UserID] = stat
	}
	stat.Count++
	stat.Value += event.Value
	updated := *stat
	sh.mu.Unlock()

	msg, err := json.Marshal(struct {
		Type  string   `json:"type"`
		Event Event    `json:"event"`
		Stat  UserStat `json:"stat"`
	}{
		Type:  "event_applied",
		Event: event,
		Stat:  updated,
	})
	if err == nil {
		s.streams.publish(msg)
	}
}

func shardIndex(userID string, shardCount int) int {
	var sum uint32
	for i := 0; i < len(userID); i++ {
		sum = sum*33 + uint32(userID[i])
	}
	return int(sum % uint32(shardCount))
}

func (s *server) ingestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if event.UserID == "" || event.Action == "" {
		http.Error(w, "user_id and action are required", http.StatusBadRequest)
		return
	}
	if event.Received.IsZero() {
		event.Received = time.Now().UTC()
	}

	select {
	case s.queue <- event:
		writeJSON(w, http.StatusAccepted, map[string]any{
			"accepted": true,
			"queued":   len(s.queue),
		})
	default:
		s.dropped.Add(1)
		http.Error(w, "queue full", http.StatusServiceUnavailable)
	}
}

func (s *server) batchIngestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var events []Event
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	var accepted int
	for i := range events {
		if events[i].UserID == "" || events[i].Action == "" {
			continue
		}
		if events[i].Received.IsZero() {
			events[i].Received = time.Now().UTC()
		}

		select {
		case s.queue <- events[i]:
			accepted++
		default:
			s.dropped.Add(1)
		}
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted": accepted,
		"received": len(events),
		"queued":   len(s.queue),
	})
}

func (s *server) statsHandler(w http.ResponseWriter, _ *http.Request) {
	topN := 10
	users := s.topUsers(topN)

	writeJSON(w, http.StatusOK, snapshot{
		Processed: s.processed.Load(),
		Dropped:   s.dropped.Load(),
		Queued:    len(s.queue),
		Workers:   runtime.NumCPU(),
		Uptime:    time.Since(s.startedAt).String(),
		TopUsers:  users,
	})
}

func (s *server) leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	writeJSON(w, http.StatusOK, s.topUsers(limit))
}

func (s *server) topUsers(limit int) []UserStat {
	all := make([]UserStat, 0, 128)
	for i := range s.shards {
		s.shards[i].mu.RLock()
		for _, stat := range s.shards[i].data {
			all = append(all, *stat)
		}
		s.shards[i].mu.RUnlock()
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Value == all[j].Value {
			return all[i].Count > all[j].Count
		}
		return all[i].Value > all[j].Value
	})

	if len(all) > limit {
		all = all[:limit]
	}
	return all
}

func (s *server) streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := s.streams.subscribe()
	defer s.streams.unsubscribe(client)

	fmt.Fprint(w, "event: ready\ndata: connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-client:
			fmt.Fprintf(w, "event: update\ndata: %s\n\n", sanitizeSSEData(msg))
			flusher.Flush()
		}
	}
}

func sanitizeSSEData(msg []byte) string {
	return strings.ReplaceAll(string(msg), "\n", "")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	workerCount := runtime.NumCPU() * 2
	queueSize := 100_000
	shardCount := 256

	svc := newServer(workerCount, queueSize, shardCount)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "High-throughput event service")
		fmt.Fprintln(w, "POST /ingest")
		fmt.Fprintln(w, "POST /ingest/batch")
		fmt.Fprintln(w, "GET  /stats")
		fmt.Fprintln(w, "GET  /leaderboard?limit=20")
		fmt.Fprintln(w, "GET  /stream")
	})
	mux.HandleFunc("/ingest", svc.ingestHandler)
	mux.HandleFunc("/ingest/batch", svc.batchIngestHandler)
	mux.HandleFunc("/stats", svc.statsHandler)
	mux.HandleFunc("/leaderboard", svc.leaderboardHandler)
	mux.HandleFunc("/stream", svc.streamHandler)

	httpServer := &http.Server{
		Addr:              ":8002",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		log.Printf("listening on %s with %d workers", httpServer.Addr, workerCount)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	close(svc.queue)
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
