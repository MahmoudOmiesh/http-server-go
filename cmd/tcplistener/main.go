package main

import (
	"fmt"
	"http-server/internal/request"
	"log"
	"net"
)

const port = ":42069"

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

		reqeust, err := request.RequestFromReader(con)

		if err != nil {
			log.Fatal("failed to read request", err)
		}

		fmt.Println("Request Line:")
		fmt.Printf("- Method: %s\n", reqeust.RequestLine.Method)
		fmt.Printf("- Target: %s\n", reqeust.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", reqeust.RequestLine.HttpVersion)

		// fmt.Println("Headers:")
		// for key, value := range reqeust.Headers.headers {
		// 	fmt.Printf("- %s: %s\n", key, value)
		// }
	}

}
