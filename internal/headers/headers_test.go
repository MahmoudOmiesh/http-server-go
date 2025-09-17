package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ := headers.Get("host")
	assert.Equal(t, "localhost:42069", host)

	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid multiple value header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFoo: bar\r\nFoo: foobar\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ = headers.Get("host")
	assert.Equal(t, "localhost:42069", host)

	foo, _ := headers.Get("foo")
	assert.Equal(t, "bar,foobar", foo)
	assert.Equal(t, 48, n)
	assert.True(t, done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFoo:     bar    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ = headers.Get("host")
	assert.Equal(t, "localhost:42069", host)

	foo, _ = headers.Get("foo")
	assert.Equal(t, "bar", foo)
	assert.Equal(t, 43, n)
	assert.True(t, done)

	// Test: invalid field name token
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
