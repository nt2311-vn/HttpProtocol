package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/nt2311-vn/HttpProtocol/internal/request"
	"github.com/nt2311-vn/HttpProtocol/internal/response"
)

type Server struct {
	handler Handler
	Listen  net.Listener
	closed  atomic.Bool
}

type HandlerError struct {
	Message    string
	StatusCode response.StatusCode
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)

	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	s := &Server{handler: handler, Listen: l}
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.Listen.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.Listen.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("Errror on accept connection: ", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		heErr := &HandlerError{StatusCode: response.StatusCodeBadRequest, Message: err.Error()}
		heErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	heErr := s.handler(buf, req)

	if heErr != nil {
		heErr.Write(conn)
		return
	}

	b := buf.Bytes()

	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
}
