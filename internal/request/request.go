package request

import (
	"bytes"
	"errors"
	"http-server/internal/headers"
	"io"
	"unicode"
)

const bufferSize = 8

var SEPARATOR = "\r\n"
var ErrorInvalidRequest = errors.New("request is invalid")
var ErrorInvalidRequestLine = errors.New("request line is invalid")
var ErrorInvalidHttpVersion = errors.New("http version is not supported")
var ErrorInvalidRequestTarget = errors.New("request target is invalid")
var ErrorInvalidMethod = errors.New("method is invalid")

type RequestState string

const (
	RequestStateInit    RequestState = "initialized"
	ReqeustStateHeaders RequestState = "headers"
	RequestStateDone    RequestState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       RequestState
}

func (r *RequestLine) isValidHttpVersion() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func (r *RequestLine) isValidRequestTarget() bool {
	for _, r := range r.RequestTarget {
		if unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

func (r *RequestLine) isValidMethod() bool {
	for _, r := range r.Method {
		if unicode.IsLower(r) {
			return false
		}
	}

	return true
}

func (r *RequestLine) isValid() (bool, error) {
	if !r.isValidHttpVersion() {
		return false, ErrorInvalidHttpVersion
	}
	if !r.isValidRequestTarget() {
		return false, ErrorInvalidRequestTarget
	}
	if !r.isValidMethod() {
		return false, ErrorInvalidMethod
	}

	return true, nil
}

func (r *Request) parse(data []byte) (int, error) {
	readBytes := 0

outer:
	for {
		currentData := data[readBytes:]

		if len(currentData) == 0 {
			break outer
		}

		switch r.state {
		case RequestStateInit:
			requestLine, bytesConsumed, err := parseRequestLine(currentData)

			if err != nil {
				return 0, err
			}

			if bytesConsumed == 0 {
				break outer
			}

			r.RequestLine = *requestLine
			r.state = ReqeustStateHeaders
			readBytes += bytesConsumed

		case ReqeustStateHeaders:
			bytesConsumed, done, err := r.Headers.Parse(currentData)

			if err != nil {
				return 0, err
			}

			if bytesConsumed == 0 {
				break outer
			}

			readBytes += bytesConsumed

			if done {
				r.state = RequestStateDone
			}

		case RequestStateDone:
			break outer

		default:
			panic("shouldn't really happend")
		}
	}

	return readBytes, nil
}

func (r *Request) done() bool {
	return r.state == RequestStateDone
}

func newRequest() *Request {
	return &Request{
		state:   RequestStateInit,
		Headers: headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := newRequest()

	for !request.done() {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])

		if err != nil {
			if err == io.EOF {
				request.state = RequestStateDone
				break
			}

			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buf[:readToIndex])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:readToIndex])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	separatorIndex := bytes.Index(request, []byte(SEPARATOR))

	if separatorIndex == -1 {
		// still need more data
		return nil, 0, nil
	}

	requestLineBytes := request[:separatorIndex]
	requestLineParts := bytes.Split(requestLineBytes, []byte(" "))
	readBytes := separatorIndex + len(SEPARATOR)

	if len(requestLineParts) != 3 {
		return nil, 0, ErrorInvalidRequestLine
	}

	requestLine := RequestLine{
		Method:        string(requestLineParts[0]),
		RequestTarget: string(requestLineParts[1]),
		HttpVersion:   string(requestLineParts[2]),
	}

	isValid, err := requestLine.isValid()

	if !isValid {
		return nil, 0, err
	}

	return &requestLine, readBytes, nil
}
