package server

import (
	"fmt"
	"http-server/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func newServer(l net.Listener) *Server {
	return &Server{
		listener: l,
		isClosed: atomic.Bool{},
	}
}

func Serve(port uint16) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := newServer(l)

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()

	if err != nil {
		return err
	}

	s.isClosed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if s.isClosed.Load() {
			return
		}

		if err != nil {
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	err := response.WriteStatusLine(conn, response.StatusOk)

	if err != nil {
		return
	}

	defaultHeaders := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, defaultHeaders)

	if err != nil {
		return
	}
}
