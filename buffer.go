package riff

import (
	"sync"
)

const bufferSize = 1024

// Buffers are pooled to reduce allocations.
var bufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, bufferSize)
		return &buf
	},
}

func getBuffer() *[]byte {
	b, _ := bufferPool.Get().(*[]byte)
	*b = (*b)[:0] // Reset the underlying slice
	return b
}

func putBuffer(buf *[]byte) {
	bufferPool.Put(buf)
}
