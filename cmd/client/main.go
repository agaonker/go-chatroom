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

	cli := client.NewClient()
	if err := cli.Connect(*host, *port); err != nil {
		log.Fatal(err)
	}
	cli.Start()
}
