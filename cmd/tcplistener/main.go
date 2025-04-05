package main

import (
	"fmt"
	"log"
	"net"

	"github.com/clement-bramy/httcp/internal/request"
)

const (
	MESSAGES_FILE = "messages.txt"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("failed to create tcp listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept incoming network connection: %v", err)
		}
		defer conn.Close()
		fmt.Println("accepted new connection..")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("failed to parse incoming data: %v\n", err)
			continue
		}
		fmt.Printf(
			"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)
		fmt.Println("connection has been closed..")
	}
}
