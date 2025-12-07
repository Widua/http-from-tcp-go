package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/widua/http-from-tcp-go/internal/request"
	"github.com/widua/http-from-tcp-go/internal/response"
)

type Server struct {
	Port     int
	Listener net.Listener
	Handler  Handler
	opened   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))

	if err != nil {
		return nil, err
	}
	server := Server{Port: port, Listener: listener, Handler: handler}
	server.opened.Store(true)
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.opened.Store(false)
	return s.Listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if !s.opened.Load() {
				return
			}
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	writer := response.NewWriter(conn)
	if err != nil {
		error := HandlerError{response.BAD_REQUEST, fmt.Sprintf("Error processing request: %v", err.Error())}
		error.Write(writer)
	}

	s.Handler(writer, req)
}
