package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		defer close(linesCh)
		defer f.Close()

		line := ""
		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)

			if err != nil {
				break
			}

			buf = buf[:n]
			if i := bytes.IndexByte(buf, '\n'); i != -1 {
				line += string(buf[:i])
				buf = buf[i+1:]
				linesCh <- line
				line = ""
			}

			line += string(buf)
		}

		if line != "" {
			linesCh <- line
		}
	}()

	return linesCh
}

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	linesCh := getLinesChannel(file)

	for line := range linesCh {
		fmt.Printf("read: %s\n", line)
	}
}
