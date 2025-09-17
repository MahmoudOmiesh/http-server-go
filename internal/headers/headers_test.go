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

	assert.Equal(t, "localhost:42069", headers["host"])

	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid multiple value header
	// headers = NewHeaders()
	// data = []byte("Host: localhost:42069\r\nFoo: bar\r\nFoo: foobar\r\n\r\n")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)

	// assert.Equal(t, "localhost:42069", headers["Host"])

	// assert.Equal(t, "bar,foobar", headers["Foo"])

	// assert.Equal(t, 48, n)
	// assert.True(t, done)

	// Test: Valid 2 headers
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFoo:     bar    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "bar", headers["foo"])
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
