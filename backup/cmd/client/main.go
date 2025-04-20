package main

import (
	"flag"
	"log"

	"github.com/ashish/go-chatroom/internal/client"
)

func main() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	chatClient := client.NewClient()

	// Connect to the server
	if err := chatClient.Connect(*host, *port); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	// Start the client
	chatClient.Start()
}
