package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/clement-bramy/httcp/internal/headers"
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

		writeReqLine(os.Stdout, req.RequestLine)
		writeHeaders(os.Stdout, req.Headers)
		writeBody(os.Stdout, req.Body)

		fmt.Println("connection has been closed..")
	}
}

func writeReqLine(w io.Writer, rl request.RequestLine) {
	fmt.Fprintf(
		w,
		"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
		rl.Method,
		rl.RequestTarget,
		rl.HttpVersion,
	)
}

func writeHeaders(w io.Writer, h headers.Headers) {
	fmt.Fprintln(w, "Headers:")
	for key, value := range h {
		fmt.Fprintf(w, "- %s: %s\n", key, value)
	}

}

func writeBody(w io.Writer, body []byte) {
	fmt.Fprintln(w, "Body:")

	if len(body) > 0 {
		fmt.Fprintf(w, "%s\n", string(body))
	}
}
