package response

import (
	"fmt"
	"io"

	"github.com/widua/http-from-tcp-go/internal/headers"
)

type writerState int

const (
	INITIALIZED writerState = iota
	WRITING_HEADERS
	WRITING_BODY
	CHUNKED_DONE
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
	err := w.headerStyleWriting(headers)
	w.state = WRITING_BODY
	return err

}

func (w *Writer) headerStyleWriting(headers headers.Headers) error {
	for k, v := range headers {
		headerLine := fmt.Sprintf("%v: %v\r\n", k, v)
		_, err := w.Writer.Write([]byte(headerLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))

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
	if w.state != WRITING_BODY {
		return 0, fmt.Errorf("Chunked body we wrote after headers")
	}
	lenPart := fmt.Sprintf("%x\r\n", len(data))
	data = append(data, []byte("\r\n")...)

	nl, err := w.Writer.Write([]byte(lenPart))
	if err != nil {
		return 0, err
	}
	nb, err := w.Writer.Write(data)
	if err != nil {
		return nl, err
	}

	return nl + nb, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.state = CHUNKED_DONE
	return w.Writer.Write([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(headers headers.Headers) error {
	if w.state != CHUNKED_DONE {
		return fmt.Errorf("Trailers go after Chunked body")
	}
	err := w.headerStyleWriting(headers)

	return err
}
