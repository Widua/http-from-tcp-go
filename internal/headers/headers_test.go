package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Done after parse
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Malformed header
	headers = NewHeaders()
	data = []byte("host localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid multiple headers
	headers = NewHeaders()
	data = []byte("Host: localhost\r\nPort:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost", headers["host"])
	assert.Equal(t, 17, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "42069", headers["port"])
	assert.Equal(t, 12, n)
	assert.False(t, done)

	// Test: Invalid character in field-name
	headers = NewHeaders()
	data = []byte("H@st: localhost\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple values (RFC 9110 5.2)
	headers = NewHeaders()
	data = []byte("Host: localhost\r\nHost:127.0.0.1\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost", headers["host"])
	_, done, err = headers.Parse(data[n:])
	require.NoError(t, err)
	assert.Equal(t, "localhost, 127.0.0.1", headers["host"])
}

func TestHeaderGetter(t *testing.T) {

	// Test: Correct get
	headers := NewHeaders()
	headers["host"] = "localhost"
	value, err := headers.Get("host")
	require.NoError(t, err)
	assert.Equal(t, "localhost", value)

	// Test: Correct get uppercase
	headers = NewHeaders()
	headers["host"] = "localhost"
	value, err = headers.Get("Host")
	require.NoError(t, err)
	assert.Equal(t, "localhost", value)

	// Test: Error incorrect name
	headers = NewHeaders()
	headers["host"] = "localhost"
	value, err = headers.Get("h@st")
	require.Error(t, err)

	// Test: Error not existing header
	headers = NewHeaders()
	value, err = headers.Get("host")
	require.Error(t, err)
}
