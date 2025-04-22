# Go Chat Room

A simple terminal-based chat room application written in Go. The application consists of a server and multiple clients that can communicate with each other in real-time.

## Features

- Real-time messaging between multiple clients
- Unique usernames for each client
- Graceful server shutdown
- Error handling for network issues and client disconnections
- Simple terminal-based interface

## Building

To build both the server and client applications:

```bash
# Build server
go build -o bin/server ./cmd/server

# Build client
go build -o bin/client ./cmd/client
```

## Running

### Starting the Server

```bash
./bin/server -port 8080
```

The server will start listening on the specified port (default: 8080).

### Starting a Client

```bash
./bin/client -host localhost -port 8080
```

The client will:
1. Prompt you for a username
2. Connect to the server
3. Allow you to send and receive messages

### Commands

- Type any message and press Enter to send it to all connected clients
- Type `exit` to leave the chat room
- Press Ctrl+C to gracefully shut down the client or server

## Project Structure

```
.
├── cmd/
│   ├── client/     # Client application
│   └── server/     # Server application
├── internal/
│   ├── client/     # Client package
│   └── server/     # Server package
└── README.md
```

## Dependencies

The project uses only standard Go libraries:
- `net` for TCP communication
- `os` and `fmt` for terminal I/O 

We use - `chan`  for message passing and concurrency.
