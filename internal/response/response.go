package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/widua/http-from-tcp-go/internal/headers"
)

type Status int
type writerState int

const protocol = "HTTP/1.1"
const (
	INITIALIZED writerState = iota
	WRITING_HEADERS
	WRITING_BODY
)
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

type Writer struct {
	Writer io.Writer
	state  writerState
}

func NewWriter(writer io.Writer) Writer {
	return Writer{
		Writer: writer,
		state:  INITIALIZED,
	}
}

func (w *Writer) WriteStatusLine(status Status) error {
	if w.state != INITIALIZED {
		return fmt.Errorf("Status line should be writed first")
	}
	statusCode := statusMap[status]
	statusLine := fmt.Sprintf("%v %v %v\r\n", protocol, statusCode.Status, statusCode.Reason)
	_, err := w.Writer.Write([]byte(statusLine))
	w.state = WRITING_HEADERS
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WRITING_HEADERS {
		return fmt.Errorf("Headers should be written after status line, and before body")
	}
	for k, v := range headers {
		headerLine := fmt.Sprintf("%v: %v\r\n", k, v)
		_, err := w.Writer.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	w.state = WRITING_BODY
	return err

}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != WRITING_BODY {
		return 0, fmt.Errorf("Body should be wrote last, after headers")
	}

	n, err := w.Writer.Write(body)

	return n, err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}
