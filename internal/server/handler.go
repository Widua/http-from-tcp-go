package server

import (
	"github.com/widua/http-from-tcp-go/internal/request"
	"github.com/widua/http-from-tcp-go/internal/response"
)

type Handler func(w response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.Status
	Message    string
}

func (err HandlerError) Write(w response.Writer) {
	message := []byte(err.Message)
	headers := response.GetDefaultHeaders(len(message))
	w.WriteStatusLine(err.StatusCode)
	w.WriteHeaders(headers)
	w.WriteBody(message)
}
