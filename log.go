package requests

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type (
	RequestPrinter interface {
		LogRequest(method string, url string, headers http.Header)
	}

	ResponsePrinter interface {
		LogResponse(statusCode int, headers http.Header, data []byte)
	}

	stdPrinter struct {
		logger *log.Logger
	}
)

var (
	_ RequestPrinter  = (*stdPrinter)(nil)
	_ ResponsePrinter = (*stdPrinter)(nil)

	defaultRequestPrinter  RequestPrinter
	defaultResponsePrinter ResponsePrinter
)

func newStdPrinter() *stdPrinter {
	return &stdPrinter{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (p *stdPrinter) LogRequest(method string, url string, headers http.Header) {
	p.logger.Output(2, fmt.Sprintf("method=%s, request_url=%s, header=%v", method, url, headers))
}

func (p *stdPrinter) LogResponse(statusCode int, headers http.Header, data []byte) {
	p.logger.Output(2, fmt.Sprintf("status_code=%d, header=%v, body=%s", statusCode, headers, string(data)))
}

func init() {
	printer := newStdPrinter()
	defaultRequestPrinter = printer
	defaultResponsePrinter = printer
}
