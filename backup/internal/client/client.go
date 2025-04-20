package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Client represents a chat client
type Client struct {
	conn     net.Conn
	username string
	done     chan struct{}
}

// NewClient creates a new chat client
func NewClient() *Client {
	return &Client{
		done: make(chan struct{}),
	}
}

// Connect connects to the chat server
func (c *Client) Connect(host, port string) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	c.conn = conn

	// Get username from user
	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	c.username = scanner.Text()

	// Send username to server
	_, err = c.conn.Write([]byte(c.username))
	if err != nil {
		return fmt.Errorf("failed to send username: %v", err)
	}

	return nil
}

// Start starts the client's message handling
func (c *Client) Start() {
	// Handle incoming messages
	go c.handleIncomingMessages()

	// Handle user input
	go c.handleUserInput()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	c.Shutdown()
}

// handleIncomingMessages handles messages from the server
func (c *Client) handleIncomingMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			message, err := bufio.NewReader(c.conn).ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("Server closed the connection")
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}
			fmt.Print(message)
		}
	}
}

// handleUserInput handles user input and sends messages to the server
func (c *Client) handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-c.done:
			return
		default:
			message := scanner.Text()
			_, err := c.conn.Write([]byte(message))
			if err != nil {
				log.Printf("Error sending message: %v", err)
				continue
			}
			if message == "exit" {
				return
			}
		}
	}
}

// Shutdown gracefully shuts down the client
func (c *Client) Shutdown() {
	close(c.done)
	if c.conn != nil {
		c.conn.Close()
	}
}
