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
		log.Fatal("couldn't resolve udp address", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatal("couldn't prepare connection", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')

		if err != nil {
			log.Print("couldn't read line", err)
			break
		}

		_, err = conn.Write([]byte(line))

		if err != nil {
			log.Print("couldn't send bytes", err)
			break
		}
	}
}
