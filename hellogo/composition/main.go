package main

import "fmt"

type Logger struct {
	prefix string
}

func (logger *Logger) Info(message string) {
	fmt.Printf("%s INFO: %s\n", logger.prefix, message)
}

// Server is composed with a Logger: it has a Logger dependency.
type Server struct {
	logger *Logger
}

// NewServer receives the dependency instead of creating one itself.
func NewServer(logger *Logger) *Server {
	return &Server{logger: logger}
}

func (server *Server) Start() {
	server.logger.Info("server started")
}

func main() {
	logger := &Logger{prefix: "app"}
	server := NewServer(logger)

	server.Start()
}
