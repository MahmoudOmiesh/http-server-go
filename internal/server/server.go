package server

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

	request, err := request.RequestFromReader(conn)

	if err != nil {
		handlerError := MakeHandlerError(response.StatusBadRequest, "bad request\n")
		handlerError.write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buf, request)

	if handlerError != nil {
		handlerError.write(conn)
		return
	}

	response.WriteStatusLine(conn, response.StatusOk)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, headers)
	response.WriteBody(conn, buf.Bytes())
}

func MakeHandlerError(code response.StatusCode, msg string) *HandlerError {
	return &HandlerError{
		Code:    code,
		Message: msg,
	}
}

func (h *HandlerError) write(w io.Writer) {
	response.WriteStatusLine(w, h.Code)
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	response.WriteBody(w, []byte(h.Message))
}
