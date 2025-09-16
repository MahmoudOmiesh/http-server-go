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

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("couldn't read request", err)
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
	lines := strings.Split(requestString, "\r\n")

	if len(lines) == 0 {
		return nil, errors.New("request string is invalid")
	}

	requestLineString := lines[0]
	requestLineParts := strings.Split(requestLineString, " ")

	if len(requestLineParts) != 3 {
		return nil, errors.New("request line is invalid")
	}

	method, err := parseRequestLineMethod(requestLineParts[0])
	if err != nil {
		return nil, err
	}

	target, err := parseRequestLineTarget(requestLineParts[1])
	if err != nil {
		return nil, err
	}

	version, err := parseRequestLineVersion(requestLineParts[2])
	if err != nil {
		return nil, err
	}

	requestLine := RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}

	return &requestLine, nil
}

func parseRequestLineMethod(methodString string) (string, error) {
	for _, r := range methodString {
		if !unicode.IsUpper(r) {
			return "", errors.New("method has lower case letter")
		}
	}

	return methodString, nil
}

func parseRequestLineTarget(targetString string) (string, error) {
	for _, r := range targetString {
		if unicode.IsSpace(r) {
			return "", errors.New("request target has spaces")
		}
	}

	return targetString, nil
}

func parseRequestLineVersion(versionString string) (string, error) {
	versionParts := strings.Split(versionString, "/")

	if len(versionParts) != 2 {
		return "", errors.New("invalid version string")
	}

	httpName := versionParts[0]
	httpVersion := versionParts[1]

	if httpName != "HTTP" {
		return "", errors.New("invalid http name")
	}

	if len(httpVersion) != 3 || !unicode.IsDigit(rune(httpVersion[0])) || !unicode.IsDigit(rune(httpVersion[2])) || httpVersion[1] != '.' {
		return "", errors.New("invalid http version")
	}

	return httpVersion, nil
}
