package request

import (
	"bytes"
	"errors"
	"io"
	"strconv"
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

		parsedBytes, err := req.parse(buf[:readToIndex])

		if err != nil {
			return nil, err
		}
		if parsedBytes == 0 {
			readedBytes := 0
			readedBytes, buf, err = req.readData(input, buf, readToIndex)
			readToIndex += readedBytes
		}

		newBuf := make([]byte, len(buf))
		copy(newBuf, buf[parsedBytes:])
		buf = newBuf
		readToIndex -= parsedBytes
	}

	contentLength, _ := req.getContentLength()
	if contentLength != len(req.Body) {
		return nil, errors.New("Invalid Content-Lenght")
	}
	return &req, nil
}

func (r *Request) readData(input io.Reader, buf []byte, readToIndex int) (int, []byte, error) {
	readedBytes, err := input.Read(buf[readToIndex:])
	if err == io.EOF {
		r.state = DONE
		return 0, buf, err
	}

	return readedBytes, buf, err
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
		if done || bytes.Index(data[n:], []byte(crlf)) == 0 {
			r.state = PARSING_BODY
			r.getContentLength()
			return n + len(crlf), nil
		}
		if n == 0 {
			return 0, nil
		}

		return n, nil

	case PARSING_BODY:

		contentLength, err := r.getContentLength()
		if err != nil {
			return 0, nil
		}

		r.Body = append(r.Body, data...)

		if len(r.Body) > contentLength {
			return 0, errors.New("Body greater than content-length")
		}
		if len(r.Body) == contentLength {
			r.state = DONE
			return len(data), nil
		}
		return len(data), nil

	case DONE:
		return 0, errors.New("Trying to read data in done state")

	}

	return 0, errors.New("Unknown State")

}

func (r *Request) getContentLength() (int, error) {
	contentLen, err := r.Headers.Get("content-length")
	if err != nil {
		return 0, err
	}
	if contentLen == "" {
		r.state = DONE
		return 0, nil
	}
	ln, err := strconv.Atoi(contentLen)
	if err != nil {
		return 0, err
	}
	return ln, nil
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
