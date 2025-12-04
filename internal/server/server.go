package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	Port     int
	Listener net.Listener
	opened   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))

	if err != nil {
		return nil, err
	}
	server := Server{Port: port, Listener: listener}
	server.opened.Store(true)
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.opened.Store(false)
	return s.Listener.Close()
}

func (s *Server) listen() {
	for s.opened.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			s.opened.Store(false)
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	fmt.Println("HANDLE")
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!"
	conn.Write([]byte(resp))
	conn.Close()

}
