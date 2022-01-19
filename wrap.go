package requests

import (
	"bytes"
	"io"
)

type (
	requestBodyRecorder struct {
		rw     io.ReadWriter
		buffer *bytes.Buffer
	}
)

func newRequestBodyRecorder(w io.ReadWriter) *requestBodyRecorder {
	return &requestBodyRecorder{
		rw:     w,
		buffer: GetBuffer(),
	}
}

func (r *requestBodyRecorder) Write(data []byte) (int, error) {
	n, err := r.rw.Write(data)
	r.buffer.Write(data[:n])
	return n, err
}

func (r *requestBodyRecorder) Close() error {
	PutBuffer(r.buffer)
	return nil
}

func (r *requestBodyRecorder) Read(data []byte) (int, error) {
	return r.rw.Read(data)
}

func (r *requestBodyRecorder) Dump() []byte {
	return r.buffer.Bytes()
}
