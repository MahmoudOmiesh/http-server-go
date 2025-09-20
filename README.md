# http-server

A minimal HTTP/1.1 server and parser implemented in Go, with:

- Manual request line, headers, and body parsing
- A lightweight server loop with a custom handler function
- Response writer utilities, including chunked transfer encoding and trailers
- Example programs for a simple HTML responder and a streaming proxy

This project is educational: it favors clarity and control over production-hardening.

## Features

- HTTP/1.1 request parsing:
  - Request line validation (method uppercase, no whitespace in target, version HTTP/1.1)
  - Incremental parsing with state machine: request line → headers → body
  - Body handling via `Content-Length` (required for bodies)
- Header utilities:
  - Parse line-by-line until empty line
  - Validates field-name token per RFC token charset
  - Case-insensitive keys, multi-value coalescing via comma
- Response writer:
  - Status line helpers (200/400/500 + fallback)
  - Default headers helper (length, close, content-type)
  - Write headers/body
  - Chunked encoding helpers and trailers
- Server:
  - TCP listener accept loop with per-connection goroutine
  - Minimal handler signature: `func(w *response.Writer, req *request.Request)`
  - Graceful close support
- Examples:
  - Basic HTML responder
  - Streaming proxy to `httpbin.org/stream/{n}` with chunked transfer and trailers
- Tests:
  - Request parsing across chunk boundaries
  - Header parsing incl. invalid cases
  - Body parsing with content-length checks

## Project layout

- `cmd/httpserver`: runnable server with example handlers (basic or proxy)
- `cmd/tcplistener`: debug tool that prints parsed request details over TCP
- `cmd/udpsender`: simple UDP sender (utility for experiments)
- `internal/request`: HTTP request parser and types
- `internal/headers`: header map and parser
- `internal/response`: response writer utilities
- `internal/server`: TCP server and handler integration

## Getting started

Requirements:

- Go (module sets `go 1.24.6`; Go 1.21+ should work, but use a recent Go if possible)

Build all:

```bash
go build ./...
```

Run tests:

```bash
go test ./...
```

## Running the server

The main entrypoint is `cmd/httpserver`.

- By default it starts the streaming proxy example on port `42069`.
- To run it:

```bash
go run ./cmd/httpserver
```

You’ll see:

- Server listens on `:42069`
- Graceful shutdown on SIGINT/SIGTERM

### Switching between examples

Open `cmd/httpserver/main.go`:

- To run the simple HTML responder, use `basicServer()`
- To run the streaming proxy (default), use `proxyServer()`

Example (top of `main`):

```go
// server, err := basicServer()
server, err := proxyServer()
```

## Usage examples

### 1) Streaming proxy (default)

Proxy path: `/httpbin/stream/{n}`

- Starts a GET to `https://httpbin.org/stream/{n}`
- Streams the response to the client using chunked transfer
- Sends trailers: `X-Content-SHA256` and `X-Content-Length`

Try:

```bash
curl -i --http1.1 "http://localhost:42069/httpbin/stream/5"
```

Notes:

- Response uses `Transfer-Encoding: chunked`
- Trailers are sent after the body; some clients do not display trailers by default

### 2) Basic HTML responder

Returns simple HTML based on request target. Toggle `basicServer()` as above and run:

```bash
go run ./cmd/httpserver
```

Then:

```bash
curl -i http://localhost:42069/
curl -i http://localhost:42069/yourproblem
curl -i http://localhost:42069/myproblem
```

## Debug tools

- `cmd/tcplistener`: accepts TCP connections on `:42069`, parses an HTTP request, and prints:
  - Request line (method, target, version)
  - Headers (key/value)
  - Body (as text)

```bash
go run ./cmd/tcplistener
```

- `cmd/udpsender`: simple REPL that sends lines over UDP to `localhost:42069` (not HTTP)

```bash
go run ./cmd/udpsender
```

## API overview

### Server

- `internal/server`
  - `Serve(port uint16, h Handler) (*Server, error)`
  - `type Handler func(w *response.Writer, req *request.Request)`
  - `(*Server).Close() error`

Handler error helper (used to write error responses):

- `MakeHandlerError(code response.StatusCode, msg string)`

### Request

- `internal/request`
  - `type Request struct { RequestLine; Headers; Body }`
  - `type RequestLine { Method, RequestTarget, HttpVersion }`
  - `RequestFromReader(io.Reader) (*Request, error)` — incremental parse loop
  - Validates: HTTP/1.1 only, uppercase method, no whitespace in target
  - Body requires `Content-Length` and reads exactly that many bytes

### Headers

- `internal/headers`
  - `type Headers`
  - `NewHeaders() Headers`
  - `(*Headers).Parse([]byte) (read int, done bool, err error)` — reads until empty line
  - `Get`, `Set` (coalesces dup keys with comma), `Replace`, `Delete`, `ForEach`

### Response

- `internal/response`
  - `type Writer`
  - `NewWriter(io.Writer) *Writer`
  - `WriteStatusLine(code StatusCode) error`
  - `GetDefaultHeaders(contentLen int) headers.Headers`
  - `WriteHeaders(headers.Headers) error`
  - `WriteBody([]byte) (int, error)`
  - Chunked helpers: `WriteChunkedBody`, `WriteChunkedBodyDone(hasTrailers bool)`, `WriteTrailers(headers.Headers)`

## Limitations

- HTTP/1.1 only (no HTTP/2)
- No TLS
- No keep-alive by default (`Connection: close` in default headers)
- No routing/middleware (single handler function)
- Minimal error reporting and resilience (educational code)
