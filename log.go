package requests

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type (
	RequestPrinter interface {
		LogRequest(*http.Request, []byte)
	}

	ResponsePrinter interface {
		LogResponse(*http.Response, []byte)
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

func (p *stdPrinter) LogRequest(req *http.Request, body []byte) {
	if len(body) > 1024 {
		p.logger.Output(2, fmt.Sprintf("[method]=%s, [request_url]=%s, [query]=%s, [header]=%v",
			req.Method, req.URL.String(), req.URL.RawQuery, req.Header))
		return
	}
	p.logger.Output(2, fmt.Sprintf("[method]=%s, [request_url]=%s, [query]=%s, [header]=%v, [body]=%s",
		req.Method, req.URL.String(), req.URL.RawQuery, req.Header, body))
}

func (p *stdPrinter) LogResponse(res *http.Response, body []byte) {
	if len(body) > 1024 {
		p.logger.Output(2, fmt.Sprintf("[status_code]=%d, [header]=%v", res.StatusCode, res.Header))
		return
	}
	p.logger.Output(2, fmt.Sprintf("[status_code]=%d, [header]=%v, [body]=%s", res.StatusCode, res.Header, body))
}

func init() {
	printer := newStdPrinter()
	defaultRequestPrinter = printer
	defaultResponsePrinter = printer
}
