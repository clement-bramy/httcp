package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/clement-bramy/httcp/internal/response"
)

type Server struct {
	Port     int
	listener net.Listener
	closed   *atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, fmt.Errorf("failed to start listening: %v", err)
	}

	var closed atomic.Bool

	server := &Server{
		Port:     port,
		listener: listener,
		closed:   &closed,
	}

	go server.listen()

	return server, nil
}

func (s *Server) listen() {
	fmt.Printf("started listening on port: %d\n", s.Port)
	for s.closed.Load() == false {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection (closed=%t): %v", s.closed.Load(), err)
			continue
		}

		fmt.Println("accepted new connection")
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handling new request")

	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))

	fmt.Println("done handling new request")
}

func (s *Server) Close() error {
	fmt.Println("stopping server..")
	s.closed.Store(true)
	return s.listener.Close()
}
