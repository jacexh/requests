package requests

import (
	"bytes"
	"sync"
)

type (
	bufferPool struct {
		pool *sync.Pool
	}
)

var (
	defaultBufferPool *bufferPool
)

func newBufferPool() *bufferPool {
	return &bufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

func (bp *bufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

func (bp *bufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.pool.Put(b)
}

func GetBuffer() *bytes.Buffer {
	return defaultBufferPool.Get()
}

func PutBuffer(b *bytes.Buffer) {
	defaultBufferPool.Put(b)
}

func init() {
	defaultBufferPool = newBufferPool()
}
