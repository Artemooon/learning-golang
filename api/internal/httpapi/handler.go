// Package httpapi exposes the records store as a JSON HTTP API.
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/artem/vintage-jazz-records-api/internal/records"
)

type handler struct {
	store *records.Store
}

// NewHandler creates routes for the records API.
func NewHandler(store *records.Store) http.Handler {
	h := handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /records", h.listRecords)
	mux.HandleFunc("POST /records", h.createRecord)
	mux.HandleFunc("GET /records/{id}", h.getRecord)
	return mux
}

func (h handler) listRecords(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.All())
}

func (h handler) createRecord(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var record records.Record
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&record); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return
	}
	if err := validate(record); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.store.Add(record); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Set("Location", "/records/"+record.ID)
	writeJSON(w, http.StatusCreated, record)
}

func (h handler) getRecord(w http.ResponseWriter, r *http.Request) {
	record, err := h.store.ByID(r.PathValue("id"))
	if errors.Is(err, records.ErrNotFound) {
		writeError(w, http.StatusNotFound, "record not found")
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func validate(record records.Record) error {
	if strings.TrimSpace(record.ID) == "" || strings.TrimSpace(record.Title) == "" || strings.TrimSpace(record.Artist) == "" || record.Year < 1900 || record.Year > 1999 {
		return errors.New("id, title, artist, and a year from 1900 through 1999 are required")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
