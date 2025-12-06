package server

import (
	"io"

	"github.com/widua/http-from-tcp-go/internal/request"
	"github.com/widua/http-from-tcp-go/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.Status
	Message    string
}

func (err HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, err.StatusCode)
	message := []byte(err.Message)
	headers := response.GetDefaultHeaders(len(message))
	response.WriteHeaders(w, headers)
	w.Write(message)
}
