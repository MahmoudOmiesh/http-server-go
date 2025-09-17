package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const SEPARATOR = "\r\n"

var (
	ErrorInvalidFieldLine = errors.New("invalid field line")
	ErrorInvalidFieldName = errors.New("invalid field name")
)

func parseFieldName(fieldName []byte) (string, error) {
	re := regexp.MustCompile("^[A-Za-z0-9!#$%&'*+\\-.\\^_`|~]+$")
	fieldNameStr := string(fieldName)

	if !re.MatchString(fieldNameStr) {
		return "", ErrorInvalidFieldName
	}

	return strings.ToLower(fieldNameStr), nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	fieldLineParts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(fieldLineParts) != 2 {
		return "", "", ErrorInvalidFieldLine
	}

	fieldName, err := parseFieldName(fieldLineParts[0])

	if err != nil {
		return "", "", err
	}

	fieldValue := bytes.TrimSpace(fieldLineParts[1])

	return string(fieldName), string(fieldValue), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	readBytes := 0
	done := false

	for {
		separatorIndex := bytes.Index(data[readBytes:], []byte(SEPARATOR))

		if separatorIndex == -1 {
			return readBytes, false, nil
		}

		// empty header
		if separatorIndex == 0 {
			done = true
			readBytes += len(SEPARATOR)
			break
		}

		fieldName, fieldValue, err := parseHeader(data[readBytes : readBytes+separatorIndex])

		if err != nil {
			return 0, false, err
		}

		h[fieldName] = fieldValue
		readBytes += separatorIndex + len(SEPARATOR)
	}

	return readBytes, done, nil
}
