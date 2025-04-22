package server

import (
	"fmt"
	"log"
	"net"
)

// Message represents a chat message
type Message struct {
	text   string
	sender *Client
}

// Client represents a connected client
type Client struct {
	conn     net.Conn
	username string
	messages chan string
}

// Server represents the chat server
type Server struct {
	clients    map[*Client]bool
	messages   chan Message
	newClients chan *Client
	delClients chan *Client
}

// NewServer creates a new chat server
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		messages:   make(chan Message),
		newClients: make(chan *Client),
		delClients: make(chan *Client),
	}
}

// Start starts the server on the specified port
func (s *Server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	log.Printf("Server started on port %s", port)

	// Start the message broadcaster
	go s.broadcastMessages()

	// Accept connections
	go s.acceptConnections(listener)

	return nil
}

// acceptConnections accepts new client connections
func (s *Server) acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		client := &Client{
			conn:     conn,
			messages: make(chan string),
		}
		s.newClients <- client
		go s.handleClient(client)
	}
}

// handleClient handles communication with a single client
func (s *Server) handleClient(client *Client) {
	// Get username from client
	buf := make([]byte, 1024)
	n, err := client.conn.Read(buf)
	if err != nil {
		log.Printf("Error reading username: %v", err)
		s.delClients <- client
		return
	}
	client.username = string(buf[:n])

	// Broadcast join message
	s.messages <- Message{
		text:   fmt.Sprintf("%s has joined the chat", client.username),
		sender: client,
	}

	// Start client message writer
	go s.writeClientMessages(client)

	// Handle messages from client
	for {
		n, err := client.conn.Read(buf)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			s.delClients <- client
			return
		}

		message := string(buf[:n])
		if message == "exit" {
			s.messages <- Message{
				text:   fmt.Sprintf("%s has left the chat", client.username),
				sender: client,
			}
			s.delClients <- client
			return
		}

		s.messages <- Message{
			text:   fmt.Sprintf("%s: %s", client.username, message),
			sender: client,
		}
	}
}

// writeClientMessages writes messages to a client
func (s *Server) writeClientMessages(client *Client) {
	for msg := range client.messages {
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Printf("Error writing to client: %v", err)
			s.delClients <- client
			return
		}
	}
}

// broadcastMessages handles message broadcasting
func (s *Server) broadcastMessages() {
	for {
		select {
		case msg := <-s.messages:
			// Broadcast message to all clients except sender
			for client := range s.clients {
				if client != msg.sender {
					client.messages <- msg.text
				}
			}

		case client := <-s.newClients:
			s.clients[client] = true

		case client := <-s.delClients:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.messages)
				client.conn.Close()
			}
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	for client := range s.clients {
		client.conn.Close()
		close(client.messages)
	}
}
