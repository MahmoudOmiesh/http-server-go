package main

import (
	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	// server, err := basicServer()
	server, err := proxyServer()

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func basicServer() (*server.Server, error) {
	return server.Serve(port, func(w *response.Writer, req *request.Request) {
		var msg []byte
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			msg = respone400()
		case "/myproblem":
			msg = respone500()
		default:
			msg = respone200()
		}

		w.WriteStatusLine(response.StatusBadRequest)
		headers := response.GetDefaultHeaders(len(msg))
		headers.Replace("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(msg)
	})
}

func proxyServer() (*server.Server, error) {
	return server.Serve(port, func(w *response.Writer, req *request.Request) {
		route := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		if route == req.RequestLine.RequestTarget {
			log.Print("request isn't to httpbin servers")
			return
		}

		res, err := http.Get("https://httpbin.org/" + route)

		if err != nil {
			log.Print("something went wrong while getting data", err)
			return
		}

		buf := make([]byte, 1024)
		defer res.Body.Close()

		w.WriteStatusLine(response.StatusOk)

		headers := response.GetDefaultHeaders(0)
		headers.Delete("Content-Length")
		headers.Set("Transfer-Encoding", "chunked")
		w.WriteHeaders(headers)

		for {
			n, err := res.Body.Read(buf)

			if n > 0 {
				w.WriteChunkedBody(buf[:n])
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Print("something went wrong while reading data", err)
				return
			}
		}

		w.WriteChunkedBodyDone()
	})
}

func respone200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func respone400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respone500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
