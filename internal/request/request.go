package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type parseState int

const crlf = "\r\n"

const (
	INITIALIZED parseState = iota
	DONE
)

type Request struct {
	RequestLine RequestLine
	State       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(input io.Reader) (*Request, error) {
	fullData, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}
	requestLine, err := parseRequestLine(fullData)

	if err != nil {
		return nil, fmt.Errorf("Eror while parsing request line: %v", err)
	}
	req := Request{RequestLine: *requestLine}
	return &req, nil
}

func parseRequestLine(reqData []byte) (*RequestLine, error) {

	ix := bytes.Index(reqData, []byte(crlf))

	if ix == -1 {
		return nil, errors.New("Could not find CRLF in request line")
	}

	requestline := string(reqData[:ix])

	return requestLineFromString(requestline)
}

func requestLineFromString(requestLine string) (*RequestLine, error) {
	unsupportedHTTPVersion := errors.New("Unsupported. Only HTTP version 1.1")
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
