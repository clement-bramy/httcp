package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("failed to resolve udp address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to dial up [%v]: %v", addr, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("failed to read input: %v\n", err)
			continue
		}

		_, err = conn.Write([]byte(in))
		if err != nil {
			log.Printf("failed to send udp content: %v", err)
			continue
		}

	}
}
