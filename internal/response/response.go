package response

import (
	"github.com/widua/http-from-tcp-go/internal/headers"
	"strconv"
)

type Status int

const protocol = "HTTP/1.1"

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

type StatusCode struct {
	Status Status
	Reason string
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}

func GetChunkEncodingHeaders() headers.Headers {
	headers := headers.NewHeaders()
	headers["connection"] = "close"
	headers["transfer-encoding"] = "chunked"
	headers["content-type"] = "text/plain"
	return headers
}
