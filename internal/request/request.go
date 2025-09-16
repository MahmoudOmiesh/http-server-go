package request

import (
	"errors"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
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

var SEPARATOR = "\r\n"
var ErrorInvalidRequest = errors.New("request is invalid")
var ErrorInvalidRequestLine = errors.New("request line is invalid")
var ErrorInvalidHttpVersion = errors.New("http version is not supported")
var ErrorInvalidRequestTarget = errors.New("request target is invalid")
var ErrorInvalidMethod = errors.New("method is invalid")

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("unable to read io.ReadAll", err)
	}

	requestString := string(requestBytes)
	requestLine, err := parseRequestLine(requestString)

	if err != nil {
		return nil, err
	}

	r := Request{RequestLine: *requestLine}

	return &r, nil
}

func parseRequestLine(requestString string) (*RequestLine, error) {
	lines := strings.Split(requestString, SEPARATOR)

	if len(lines) == 0 {
		return nil, ErrorInvalidRequest
	}

	requestLineString := lines[0]
	requestLineParts := strings.Split(requestLineString, " ")

	if len(requestLineParts) != 3 {
		return nil, ErrorInvalidRequestLine
	}

	requestLine := RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   requestLineParts[2],
	}

	isValid, err := requestLine.isValid()

	if !isValid {
		return nil, err
	}

	return &requestLine, nil
}
