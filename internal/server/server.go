package server

import (
	"bytes"
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
	if err != nil {
		error := HandlerError{response.BAD_REQUEST, fmt.Sprintf("Error processing request: %v", err.Error())}
		error.Write(conn)
	}
	buff := bytes.NewBuffer(make([]byte, 0))
	reqErr := s.Handler(buff, req)
	if reqErr != nil {
		reqErr.Write(conn)
		return
	}
	body := buff.Bytes()
	headers := response.GetDefaultHeaders(len(body))
	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
}
