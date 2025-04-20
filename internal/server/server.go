package server

import (
	"fmt"
	"log"
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
	clients  map[*Client]bool
	mutex    sync.RWMutex
	listener net.Listener
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
	s.listener = listener
	log.Printf("Server started on port %s", port)

	go s.acceptConnections()
	return nil
}

// acceptConnections accepts new client connections
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		client := &Client{conn: conn}
		s.mutex.Lock()
		s.clients[client] = true
		s.mutex.Unlock()

		go s.handleClient(client)
	}
}

// handleClient handles communication with a single client
func (s *Server) handleClient(client *Client) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, client)
		s.mutex.Unlock()
		client.conn.Close()
	}()

	// Get username from client
	buf := make([]byte, 1024)
	n, err := client.conn.Read(buf)
	if err != nil {
		log.Printf("Error reading username: %v", err)
		return
	}
	client.username = string(buf[:n])

	// Broadcast join message
	s.broadcast(fmt.Sprintf("%s has joined the chat", client.username), client)

	// Handle messages from client
	for {
		n, err := client.conn.Read(buf)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		message := string(buf[:n])
		if message == "exit" {
			s.broadcast(fmt.Sprintf("%s has left the chat", client.username), client)
			return
		}

		// Broadcast message to all clients
		s.broadcast(fmt.Sprintf("%s: %s", client.username, message), client)
	}
}

// broadcast sends a message to all clients except the sender
func (s *Server) broadcast(message string, sender *Client) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.clients {
		if client != sender {
			_, err := client.conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
			}
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.mutex.Lock()
	for client := range s.clients {
		client.conn.Close()
	}
	s.mutex.Unlock()
	s.listener.Close()
}
