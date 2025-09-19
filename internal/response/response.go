package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
)

type StatusCode int

type Writer struct {
	writer io.Writer
}

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func getStatusLine(statusCode StatusCode) []byte {
	switch statusCode {
	case StatusOk:
		return []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		return []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		return []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Appendf([]byte("HTTP/1.1"), " %d \r\n", statusCode)
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := getStatusLine(statusCode)

	_, err := w.writer.Write(statusLine)

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var foundErr error = nil

	headers.ForEach(func(key, val string) {
		if foundErr != nil {
			return
		}

		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", key, val)

		foundErr = err
	})

	_, err := fmt.Fprint(w.writer, "\r\n")

	if err != nil {
		foundErr = err
	}

	return foundErr
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	return fmt.Fprint(w.writer, string(body))
}
