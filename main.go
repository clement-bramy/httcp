package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	MESSAGES_FILE = "messages.txt"
)

func main() {
	file, err := os.Open(MESSAGES_FILE)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	buf := make([]byte, 8)
	for {
		_, err = file.Read(buf)
		if err == io.EOF {
			break
		}

		fmt.Fprintf(os.Stdout, "read: %s\n", buf)
	}
}
