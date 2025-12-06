package main

import (
	"io"
	"os"
	"os/signal"

	"log"
	"syscall"

	"github.com/widua/http-from-tcp-go/internal/request"
	"github.com/widua/http-from-tcp-go/internal/response"
	"github.com/widua/http-from-tcp-go/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handle)

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

func handle(w io.Writer, req *request.Request) *server.HandlerError {

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{StatusCode: response.BAD_REQUEST, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{StatusCode: response.INTERNAL_SERVER_ERROR, Message: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
	}
	return nil

}
