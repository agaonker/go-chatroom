package server

import (
	"net"
	"sync"
)

// Client represents a connected client
type Client struct {
	conn     net.Conn
	username string
}

// Server represents the chat server
type Server struct {
	clients map[*Client]bool
	mutex   sync.RWMutex
}

// NewServer creates a new chat server
func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
	}
}

// Start starts the server on the specified port
func (s *Server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer listener.Close()

	// TODO: Accept connections and handle clients
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	// TODO: Close all client connections
}
