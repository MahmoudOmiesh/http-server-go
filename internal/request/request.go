package request

import (
	"bytes"
	"errors"
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
	RequestStateInit RequestState = "initialized"
	RequestStateDone RequestState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
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
	if r.state == RequestStateDone {
		return 0, errors.New("parsing a request after it is done")
	}

	if r.state == RequestStateInit {
		requestLine, bytesConsumed, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if bytesConsumed == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = RequestStateDone
		return bytesConsumed, nil
	}

	return 0, errors.New("unknown request state")
}

func (r *Request) done() bool {
	return r.state == RequestStateDone
}

func newRequest() *Request {
	return &Request{
		state: RequestStateInit,
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
