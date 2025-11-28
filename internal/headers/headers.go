package headers

import (
	"bytes"
	"errors"
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

func (h Headers) Get(key string) (string, error) {
	key = strings.ToLower(key)

	if !isValidHeaderName(key) {
		return "", errors.New("Invalid header name")

	}

	val, ok := h[key]

	if !ok {
		return "", errors.New("Header not found")
	}

	return val, nil
}

func (h Headers) String() (s string) {
	s = "Headers:\n"

	for k, v := range h {
		s += fmt.Sprintf("- %v: %v\n", k, v)
	}

	return
}

func (h Headers) headerFromString(headerLine string) error {
	ix := strings.Index(headerLine, ":")

	if ix == -1 {
		return fmt.Errorf("Malformed header")
	}

	if headerLine[ix-1] == ' ' {
		return fmt.Errorf("Incorrect format of header: %v", headerLine)
	}

	headerName := strings.TrimSpace(headerLine[:ix])
	headerValue := strings.TrimSpace(headerLine[ix+1:])

	headerName = strings.ToLower(headerName)

	if !isValidHeaderName(headerName) {
		return fmt.Errorf("Invalid field-name: %v", headerName)
	}
	if v, ok := h[headerName]; ok {
		nv := fmt.Sprintf("%s, %s", v, headerValue)
		headerValue = nv
	}

	h[headerName] = headerValue
	return nil
}

func isValidHeaderName(headerName string) bool {
	const allowedChars = "!#$%&'*+-.^_`|~"

	for _, rune := range headerName {
		if rune > 127 {
			return false
		}
		if 'a' <= rune && rune <= 'z' {
			continue
		}
		if '0' <= rune && rune <= '9' {
			continue
		}
		if strings.ContainsRune(allowedChars, rune) {
			continue
		}
		return false
	}
	return true
}

func NewHeaders() Headers {
	return make(Headers)
}
