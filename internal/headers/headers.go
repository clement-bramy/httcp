package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

var (
	InvalidHeaderLineParts     = errors.New("invalid header line: parts numbers")
	InvalidKeyFormatWhitespace = errors.New("invalid header line: key has trailing whitespaces")
)

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, []byte("\r\n"))
	if index == -1 {
		return 0, false, nil
	}

	if index == 0 {
		return 2, true, nil
	}

	raw := strings.TrimSpace(string(data[:index]))
	parts := strings.SplitN(raw, ":", 2)

	key := parts[0]
	value := strings.TrimSpace(parts[1])
	if strings.HasSuffix(key, " ") {
		return 0, false, InvalidKeyFormatWhitespace
	}

	h[key] = value
	return index + 2, false, nil
}
