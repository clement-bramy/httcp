package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	MESSAGES_FILE = "messages.txt"
)

func main() {
	file, err := os.Open(MESSAGES_FILE)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	str := ""
	buf := make([]byte, 8)
	for {
		_, err = file.Read(buf)
		if err == io.EOF {
			break
		}

		str += string(buf)
		parts := strings.Split(str, "\n")
		if len(parts) > 1 {
			fmt.Printf("read: %s\n", parts[0])
			str = parts[1]
		}
	}
}
