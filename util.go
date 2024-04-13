package requests

import (
	"bytes"
	"fmt"
	"strconv"
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

func ToString(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return v.(string), nil
	case bool:
		return strconv.FormatBool(v.(bool)), nil
	case uint8:
		return strconv.FormatUint(uint64(v.(uint8)), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v.(uint16)), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v.(uint32)), 10), nil
	case uint64:
		return strconv.FormatUint(v.(uint64), 10), nil
	case uint:
		return strconv.FormatUint(uint64(v.(uint)), 10), nil
	case int:
		return strconv.FormatInt(int64(v.(int)), 10), nil
	case int8:
		return strconv.FormatInt(int64(v.(int8)), 10), nil
	case int16:
		return strconv.FormatInt(int64(v.(int16)), 10), nil
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10), nil
	case int64:
		return strconv.FormatInt(v.(int64), 10), nil
	default:
		return "", fmt.Errorf("unsupported data type: %T", t)
	}
}

func init() {
	defaultBufferPool = newBufferPool()
}
