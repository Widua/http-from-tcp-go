package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	ix := bytes.Index(data, []byte(crlf))
	if ix == -1 {
		return 0, false, nil
	}
	if ix == 0 {
		return 0, true, nil
	}
	headerLine := string(data[:ix])
	err = h.headerFromString(headerLine)
	if err != nil {
		return 0, false, err
	}

	return ix + len(crlf), false, nil
}

func (h Headers) headerFromString(headerLine string) error {
	ix := strings.Index(headerLine, ":")

	if headerLine[ix-1] == ' ' {
		return fmt.Errorf("Incorrect format of header: %v", headerLine)
	}

	headerName := strings.TrimSpace(headerLine[:ix])
	headerValue := strings.TrimSpace(headerLine[ix+1:])

	h[headerName] = headerValue
	return nil
}

func NewHeaders() Headers {
	return make(Headers)
}
