package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/nt2311-vn/HttpProtocol/internal/request"
	"github.com/nt2311-vn/HttpProtocol/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	handler Handler
	Listen  net.Listener
	closed  atomic.Bool
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
	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
}
