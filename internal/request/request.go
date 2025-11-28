package request

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/widua/http-from-tcp-go/internal/headers"
)

type parseState int

const crlf = "\r\n"
const BUFFER_SIZE = 8

const (
	INITIALIZED parseState = iota
	PARSING_HEADERS
	PARSING_BODY
	DONE
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte

	state parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(input io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE)
	readToIndex := 0
	req := Request{
		state: INITIALIZED,
	}

	for req.state != DONE {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		readedBytes, err := input.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				req.state = DONE
				break

			}
			return nil, err
		}
		readToIndex += readedBytes
		parsedBytes, err := req.parse(buf)
		if err != nil {
			return nil, err
		}
		newBuf := make([]byte, len(buf))
		copy(newBuf, buf[parsedBytes:])
		buf = newBuf

		readToIndex -= parsedBytes

	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case INITIALIZED:
		readed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if readed == 0 {
			return 0, nil
		}
		fullData := string(data[:readed])
		requestLine, err := requestLineFromString(fullData)
		if err != nil {
			return 0, err
		}
		r.RequestLine = *requestLine
		r.state = PARSING_HEADERS
		r.Headers = headers.NewHeaders()
		return readed, nil

	case PARSING_HEADERS:

		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = PARSING_BODY
			return n, nil
		}
		if n == 0 {
			return 0, nil
		}
		return n, nil

	case PARSING_BODY:

		r.state = DONE

	case DONE:
		return 0, errors.New("Trying to read data in done state")

	}

	return 0, errors.New("Unknown State")

}

func (r *Request) parseBody(data []byte) {

}

func parseRequestLine(reqData []byte) (int, error) {

	ix := bytes.Index(reqData, []byte(crlf))

	if ix == -1 {
		return 0, nil
	}

	return ix + len(crlf), nil
}

func requestLineFromString(requestLine string) (*RequestLine, error) {
	unsupportedHTTPVersion := errors.New("Unsupported. Only HTTP version 1.1")
	requestLine = strings.TrimRight(requestLine, "\r\n")
	splitted := strings.Split(requestLine, " ")

	if len(splitted) != 3 {
		return nil, errors.New("Request line should have 3 parts")
	}

	method := splitted[0]
	target := splitted[1]

	if !strings.Contains(target, "/") {
		return nil, errors.New("Target should start with /")
	}

	version := splitted[2]

	httpParts := strings.Split(version, "/")
	if httpParts[0] != "HTTP" {
		return nil, unsupportedHTTPVersion
	}
	if httpParts[1] != "1.1" {
		return nil, unsupportedHTTPVersion
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   httpParts[1],
	}, nil
}
