package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/widua/http-from-tcp-go/internal/headers"
)

type Status int

type StatusCode struct {
	Status Status
	Reason string
}

const (
	OK Status = iota
	BAD_REQUEST
	INTERNAL_SERVER_ERROR
)

var statusMap = map[Status]StatusCode{
	OK:                    {200, "OK"},
	BAD_REQUEST:           {400, "Bad Request "},
	INTERNAL_SERVER_ERROR: {500, "Internal Server Error"},
}

const protocol = "HTTP/1.1"

func WriteStatusLine(w io.Writer, status Status) error {
	statusCode := statusMap[status]

	statusLine := fmt.Sprintf("%v %v %v\r\n", protocol, statusCode.Status, statusCode.Reason)
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers {
		headerLine := fmt.Sprintf("%v: %v\r\n", k, v)
		w.Write([]byte(headerLine))
	}
	w.Write([]byte("\r\n"))
	return nil
}
