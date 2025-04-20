package client

import "net"

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
	return nil
}

// Start starts the client's message handling
func (c *Client) Start() {
	// TODO: Handle incoming messages and user input
}

// Shutdown gracefully shuts down the client
func (c *Client) Shutdown() {
	if c.conn != nil {
		c.conn.Close()
	}
}
