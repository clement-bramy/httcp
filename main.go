package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("connection has been closed..")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		str := ""
		for {
			buf := make([]byte, 8)
			l, err := f.Read(buf)
			if err == io.EOF {
				break
			}

			str += string(buf[:l])
			parts := strings.Split(str, "\n")
			if len(parts) > 1 {
				out <- parts[0]
				str = parts[1]
			}
		}
		if str != "" {
			out <- str
		}
	}()

	return out
}
