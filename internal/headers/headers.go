package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Headers map[string]string

var (
	InvalidHeaderLineParts     = errors.New("invalid header line: parts numbers")
	InvalidKeyFormatWhitespace = errors.New("invalid header line: key has trailing whitespaces")
	InvalidKeyFormatCharacters = errors.New("invalid header line: key has invalid characters")
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
	if strings.HasSuffix(key, " ") {
		return 0, false, InvalidKeyFormatWhitespace
	}

	key, err = getValidKey(parts[0])
	if err != nil {
		fmt.Printf("invalid key: [%s][%s]\n", key, string(data))
		return 0, false, InvalidKeyFormatCharacters
	}

	key = strings.ToLower(key)
	value := strings.TrimSpace(parts[1])

	if existing, found := h[key]; found {
		h[key] = strings.Join([]string{existing, value}, ", ")
	} else {
		h[key] = value
	}

	return index + 2, false, nil
}

func (h Headers) Get(key string) (value string, found bool) {
	key = strings.ToLower(key)
	value, found = h[key]
	return
}

func (h Headers) GetContentLength() (length int, found bool, err error) {
	if value, found := h["content-length"]; found {
		length, err := strconv.Atoi(value)
		if err != nil {
			return 0, false, fmt.Errorf("invalid content length: %s", value)
		}

		return length, true, nil
	}

	return 0, false, nil
}

const allowedSpecialChar = "!#$%&'*+-.^_`|~"

func getValidKey(key string) (string, error) {
	for _, r := range key {
		alpha := isAlpha(r)
		digit := isDigit(r)
		specs := strings.Contains(allowedSpecialChar, string(r))

		if !alpha && !digit && !specs {
			return "", InvalidKeyFormatCharacters
		}
	}
	return strings.ToLower(key), nil
}

func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
