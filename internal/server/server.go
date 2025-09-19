package server

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	handler  Handler
	listener net.Listener
	isClosed atomic.Bool
}

func newServer(h Handler, l net.Listener) *Server {
	return &Server{
		handler:  h,
		listener: l,
		isClosed: atomic.Bool{},
	}
}

func Serve(port uint16, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := newServer(h, l)

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

		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	request, err := request.RequestFromReader(conn)

	if err != nil {
		handlerError := MakeHandlerError(response.StatusBadRequest, err.Error())
		handlerError.write(conn)
		return
	}

	responseWriter := response.NewWriter(conn)
	s.handler(responseWriter, request)
}

func MakeHandlerError(code response.StatusCode, msg string) *HandlerError {
	return &HandlerError{
		Code:    code,
		Message: msg,
	}
}

func (h *HandlerError) write(w io.Writer) {
	responseWriter := response.NewWriter(w)

	responseWriter.WriteStatusLine(h.Code)
	headers := response.GetDefaultHeaders(len(h.Message))
	responseWriter.WriteHeaders(headers)
	responseWriter.WriteBody([]byte(h.Message))
}
