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

	srv := server.NewServer()
	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	srv.Shutdown()
	log.Println("Server shutdown complete")
}
