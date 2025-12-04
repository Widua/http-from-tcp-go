package main

import (
	"os"
	"os/signal"

	"log"
	"syscall"

	"github.com/widua/http-from-tcp-go/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
