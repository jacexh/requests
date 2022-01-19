package zaplog

import (
	"net/http"

	"github.com/jacexh/requests"
	"go.uber.org/zap"
)

type (
	ZapLogger struct {
		logger                 *zap.Logger
		requestFieldSelectors  []RequestFieldSelector
		responseFieldSelectors []ResponseFieldSelector
	}

	RequestFieldSelector  func(*http.Request, []byte) zap.Field
	ResponseFieldSelector func(*http.Response, []byte) zap.Field
)

var (
	_ requests.RequestPrinter = (*ZapLogger)(nil)
)

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: logger,
		requestFieldSelectors: []RequestFieldSelector{
			func(r *http.Request, b []byte) zap.Field {
				return zap.String("url", r.URL.String())
			},

			func(r *http.Request, b []byte) zap.Field {
				return zap.String("method", r.Method)
			},

			func(r *http.Request, b []byte) zap.Field {
				return zap.String("content-type", r.Header.Get("Content-Type"))
			},
		},
		responseFieldSelectors: []ResponseFieldSelector{
			func(r *http.Response, b []byte) zap.Field {
				return zap.Int("status_code", r.StatusCode)
			},

			func(r *http.Response, b []byte) zap.Field {
				return zap.Int64("content-length", r.ContentLength)
			},
		},
	}
}

func (zl *ZapLogger) WithRequestSelectors(selectors ...RequestFieldSelector) {
	zl.requestFieldSelectors = append(zl.requestFieldSelectors, selectors...)
}

func (zl *ZapLogger) WithResponseSelector(selectors ...ResponseFieldSelector) {
	zl.responseFieldSelectors = append(zl.responseFieldSelectors, selectors...)
}

func (zl *ZapLogger) LogRequest(req *http.Request, body []byte) {
	fields := make([]zap.Field, len(zl.requestFieldSelectors))
	for n, selector := range zl.requestFieldSelectors {
		fields[n] = selector(req, body)
	}
	zl.logger.Info("sending http request", fields...)
}

func (zl *ZapLogger) LogResponse(res *http.Response, body []byte) {
	fields := make([]zap.Field, len(zl.responseFieldSelectors))
	for n, selector := range zl.responseFieldSelectors {
		fields[n] = selector(res, body)
	}
	zl.logger.Info("received http response", fields...)
}
