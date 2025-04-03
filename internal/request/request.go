package request

import (
	"errors"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var (
	InvalidHttpHeaderFormat = errors.New("invalid HTTP header format")
	InvalidHttpRequestLine  = errors.New("invalid HTTP request line")
	InvalidHttpMethod       = errors.New("invalid HTTP method")
	InvalidHttpVersion      = errors.New("invalid HTTP version")

	validHttpMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	raw, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	parts := strings.Split(string(raw), "\r\n")
	if len(parts) < 1 {
		return &Request{}, InvalidHttpHeaderFormat
	}

	requestLine := parts[0]
	rparts := strings.Split(requestLine, " ")
	if len(rparts) != 3 {
		return &Request{}, InvalidHttpRequestLine
	}

	rl := &Request{
		RequestLine: RequestLine{
			Method:        rparts[0],
			RequestTarget: rparts[1],
			HttpVersion:   strings.TrimPrefix(rparts[2], "HTTP/"),
		},
	}

	if !slices.Contains(validHttpMethods, rl.RequestLine.Method) {
		return &Request{}, InvalidHttpMethod
	}

	if rl.RequestLine.HttpVersion != "1.1" {
		return &Request{}, InvalidHttpVersion
	}

	return rl, nil
}
