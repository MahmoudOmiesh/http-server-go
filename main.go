package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	line := ""

	for {
		buf := make([]byte, 8)
		n, err := file.Read(buf)

		if err != nil {
			break
		}

		current := string(buf[:n])
		currentSplit := strings.Split(current, "\n")

		line += currentSplit[0]

		if len(currentSplit) == 1 {
			continue
		}

		fmt.Printf("read: %s\n", line)

		for i := 1; i < len(currentSplit)-1; i++ {
			fmt.Printf("read: %s\n", currentSplit[i])
		}

		line = currentSplit[len(currentSplit)-1]
	}

	if line != "" {
		fmt.Printf("read: %s\n", line)
	}

}
