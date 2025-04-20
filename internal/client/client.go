package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Client represents a chat client
type Client struct {
	conn     net.Conn
	username string
}

// NewClient creates a new chat client
func NewClient() *Client {
	return &Client{}
}

// Connect connects to the chat server
func (c *Client) Connect(host, port string) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
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
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := c.conn.Write([]byte(message))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}
		if message == "exit" {
			break
		}
	}
}

// handleIncomingMessages handles messages from the server
func (c *Client) handleIncomingMessages() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
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

// Shutdown gracefully shuts down the client
func (c *Client) Shutdown() {
	if c.conn != nil {
		c.conn.Close()
	}
}
