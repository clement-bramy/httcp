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
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestInvalidWhitespaceHeaderKey(t *testing.T) {
	// Test: Invalid spacing header
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestValidHeaderWithWhitespaces(t *testing.T) {
	headers := NewHeaders()
	data := []byte("   Content-Type:   application/json   \r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, 40, n)
	assert.False(t, done)
}

func TestValidTwoHeaderWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	headers["Host"] = "localhost:43069"
	headers["Content-Length"] = "233"
	data := []byte("Content-Type:    application/json\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, 35, n)
	assert.False(t, done)
}

func TestValidHeadersDone(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestHeaderThenDoneValid(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-Type:application/json\r\n\r\n")

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, 31, n)
	assert.False(t, done)

	n, done, err = headers.Parse(data[n:])

	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}
