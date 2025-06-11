package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Listen net.Listener
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	s := &Server{Listen: l}
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
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!"
	conn.Write([]byte(response))
}
