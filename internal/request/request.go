package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/clement-bramy/httcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
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

	ReadAttemptDoneState = errors.New("trying to read data in a done state")
	UnknownParserState   = errors.New("error: unknown state")

	validHttpMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
)

const (
	bufferSize   = 8
	growthFactor = 2
)

type State int

const (
	Initialised State = iota
	ParseHeaders
	Done
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0
	buf := make([]byte, bufferSize, bufferSize)
	req := &Request{State: Initialised, Headers: headers.NewHeaders()}
	for req.State != Done {

		if readToIndex >= len(buf) {
			nbuf, err := reallocate(buf)
			if err != nil {
				return &Request{}, err
			}
			buf = nbuf
			// fmt.Printf("reallocated buf: [%d:%s]\n", len(buf), string(buf))
		}

		// read from readToIndex into the buffer
		read, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.State = Done
			fmt.Println("EOF has been reached!")
			break
		}

		readToIndex += read
		read, err = req.parse(buf[:readToIndex])
		if err != nil {
			return &Request{}, err
		}

		copy(buf, buf[read:])
		readToIndex -= read
	}

	return req, nil
}

func reallocate(source []byte) ([]byte, error) {
	length := len(source) * growthFactor
	dest := make([]byte, length, length)

	copied := copy(dest, source)
	if copied < len(source) {
		return []byte{}, fmt.Errorf("warn: reallocation did not copy entirely the source: [%d/%d]\n", copied, len(source))
	}

	return dest, nil
}

// State machine maintaining two states:
//   - Initialised: processing incoming data
//   - Done: read the expected data
//
// When data is *actually read*, returns the number of bytes read
func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != Done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += n
		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case Initialised:
		read, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if read == 0 {
			return 0, nil
		}

		r.RequestLine = rl
		r.State = ParseHeaders
		return read, err

	case ParseHeaders:
		read, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if read == 0 {
			return 0, nil
		}

		if done {
			r.State = Done
		}

		return read, err

	case Done:
		fmt.Printf("parser in done state - should not happen")
		return 0, ReadAttemptDoneState

	default:
		fmt.Printf("parser in unknown state - should not happen")
		return 0, UnknownParserState
	}
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	index := bytes.Index(data, []byte("\r\n"))

	// not enough data to parse the request line
	if index == -1 {
		return 0, RequestLine{}, nil
	}

	line := string(data[:index])
	rparts := strings.Split(line, " ")
	if len(rparts) != 3 {
		return 0, RequestLine{}, InvalidHttpRequestLine
	}

	rl := RequestLine{
		Method:        rparts[0],
		RequestTarget: rparts[1],
		HttpVersion:   strings.TrimPrefix(rparts[2], "HTTP/"),
	}

	if !slices.Contains(validHttpMethods, rl.Method) {
		return 0, RequestLine{}, InvalidHttpMethod
	}

	if rl.HttpVersion != "1.1" {
		return 0, RequestLine{}, InvalidHttpVersion
	}

	// data has been read/parsed returning the size consumed
	return index + 2, rl, nil
}
