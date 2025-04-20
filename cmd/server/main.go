package main

import (
	"flag"
	"log"

	"github.com/ashish/go-chatroom/internal/server"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	srv := server.NewServer()
	if err := srv.Start(*port); err != nil {
		log.Fatal(err)
	}
}
