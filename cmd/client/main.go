package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashish/go-chatroom/internal/client"
)

func main() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	cli := client.NewClient()
	if err := cli.Connect(*host, *port); err != nil {
		log.Fatal(err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start client in a goroutine
	go cli.Start()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down client...")
	cli.Shutdown()
	log.Println("Client shutdown complete")
}
