package main

import (
	"fmt"
	"github.com/widua/http-from-tcp-go/internal/request"
	"github.com/widua/http-from-tcp-go/internal/response"
	"github.com/widua/http-from-tcp-go/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func handle(w response.Writer, req *request.Request) {

	switch {
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/yourproblem"):
		handleBadRequest(w)
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/myproblem"):
		handleInternalServerError(w)
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin"):
		handleHttpBin(w, req)
	default:
		handleOk(w)

	}
}

func handleHttpBin(w response.Writer, req *request.Request) {
	httpBinPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	headers := response.GetChunkEncodingHeaders()

	resp, _ := http.Get(fmt.Sprintf("https://httpbin.org/%v", httpBinPath))
	w.WriteStatusLine(response.OK)
	w.WriteHeaders(headers)

	for {
		byteBuff := make([]byte, 1024)
		_, err := resp.Body.Read(byteBuff)
		if err != nil {
			if err == io.EOF {
				w.WriteChunkedBodyDone()
				return
			}
			continue
		}
		w.WriteChunkedBody(byteBuff)
	}

}

func handleBadRequest(w response.Writer) {

	resMessage := `
	<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>
	`

	err := w.WriteStatusLine(response.BAD_REQUEST)
	if err != nil {
		handleError(w, err)
		return
	}
	headers := response.GetDefaultHeaders(len(resMessage))
	headers["content-type"] = "text/html"
	err = w.WriteHeaders(headers)
	if err != nil {
		handleError(w, err)
		return
	}
	_, err = w.WriteBody([]byte(resMessage))
	if err != nil {
		handleError(w, err)
		return

	}
}
func handleInternalServerError(w response.Writer) {
	resMessage := `
	<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>
	`
	err := w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)
	if err != nil {
		handleError(w, err)
		return
	}
	headers := response.GetDefaultHeaders(len(resMessage))
	headers["content-type"] = "text/html"
	err = w.WriteHeaders(headers)
	if err != nil {
		handleError(w, err)
		return
	}
	_, err = w.WriteBody([]byte(resMessage))
	if err != nil {
		handleError(w, err)
		return

	}
}
func handleOk(w response.Writer) {
	resMessage := `
	<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>
	`
	err := w.WriteStatusLine(response.OK)
	if err != nil {
		handleError(w, err)
		return
	}
	headers := response.GetDefaultHeaders(len(resMessage))
	headers["content-type"] = "text/html"
	err = w.WriteHeaders(headers)
	if err != nil {
		handleError(w, err)
		return
	}
	_, err = w.WriteBody([]byte(resMessage))
	if err != nil {
		handleError(w, err)
		return

	}
}

func handleError(w response.Writer, err error) {
	handleErr := server.HandlerError{StatusCode: response.INTERNAL_SERVER_ERROR, Message: err.Error()}
	handleErr.Write(w)
}
