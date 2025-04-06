package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/clement-bramy/httcp/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error

	switch statusCode {
	case StatusOk:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalServerError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}

	return err
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()

	h["content-length"] = strconv.Itoa(contentLength)
	h["connection"] = "close"
	h["content-type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for key, val := range h {
		_, err := w.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	return nil
}
