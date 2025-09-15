package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)
	line := ""

	go func() {
		defer close(linesCh)
		defer f.Close()

		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)

			if err != nil {
				fmt.Println("it's fucked")
				break
			}

			current := string(buf[:n])
			currentSplit := strings.Split(current, "\n")

			line += currentSplit[0]

			if len(currentSplit) == 1 {
				continue
			}

			linesCh <- line

			for i := 1; i < len(currentSplit)-1; i++ {
				linesCh <- currentSplit[i]
			}

			line = currentSplit[len(currentSplit)-1]
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
