package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

const port = ":42069"

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
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("couldn't listen on port", port, err)
	}

	defer listener.Close()

	for {
		con, err := listener.Accept()

		if err != nil {
			log.Fatal("conntection refused", err)
			break
		}

		log.Print("connection accepted")

		linesCh := getLinesChannel(con)

		for line := range linesCh {
			fmt.Println(line)
		}
	}

}
