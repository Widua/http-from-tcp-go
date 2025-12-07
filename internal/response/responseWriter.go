package response

import (
	"fmt"
	"github.com/widua/http-from-tcp-go/internal/headers"
	"io"
)

type writerState int

const (
	INITIALIZED writerState = iota
	WRITING_HEADERS
	WRITING_BODY
)

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

func (w *Writer) WriteChunkedBody(data []byte) (int, error) {
	lenPart := fmt.Sprintf("%v\r\n", len(data))
	dataPart := fmt.Sprintf("%v\r\n", data)

	nl, err := w.Writer.Write([]byte(lenPart))
	if err != nil {
		return 0, err
	}
	nb, err := w.Writer.Write([]byte(dataPart))
	if err != nil {
		return nl, err
	}

	return nl + nb, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte("0\r\n\r\n"))
}
