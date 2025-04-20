package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashish/go-chatroom/internal/server"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	chatServer := server.NewServer()

	// Start the server
	if err := chatServer.Start(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	chatServer.Shutdown()
	log.Println("Server shutdown complete")
}
