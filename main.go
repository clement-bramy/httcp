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
	defer file.Close()

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		str := ""
		buf := make([]byte, 8)
		for {
			_, err := f.Read(buf)
			if err == io.EOF {
				close(out)
				break
			}

			str += string(buf)
			parts := strings.Split(str, "\n")
			if len(parts) > 1 {
				out <- parts[0]
				str = parts[1]
			}
		}
	}()

	return out
}
