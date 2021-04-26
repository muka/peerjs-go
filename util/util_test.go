package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	data := bytes.NewBuffer([]byte{})
	for data.Len() < ChunkedMTU {
		data.Write([]byte("another piece to the chunk"))
	}
	chunks := Chunk(data.Bytes())
	assert.NotEmpty(t, chunks)
}
