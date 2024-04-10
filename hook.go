package requests

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func UnmarshalJSON(v any) Unmarshaller {
	return func(data []byte) error {
		return json.Unmarshal(data, v)
	}
}

func UnmarlXML(v any) Unmarshaller {
	return func(data []byte) error {
		return xml.Unmarshal(data, v)
	}
}

func StdLogRequest(logger *log.Logger) BeforeRequestHook {
	return func(req *http.Request, body []byte) {
		if len(body) > 1024 {
			logger.Output(2, fmt.Sprintf("[method]=%s, [request_url]=%s, [header]=%v",
				req.Method, req.URL.String(), req.Header))
			return
		}
		logger.Output(2, fmt.Sprintf("[method]=%s, [request_url]=%s, [header]=%v, [body]=%s",
			req.Method, req.URL.String(), req.Header, body))
	}
}

func StdLogResponse(logger *log.Logger) AfterRequestHook {
	return func(res *http.Response, err error) {
		if err != nil {
			logger.Output(2, fmt.Sprintf("[status_code]=%d, [header]=%v, [error]=%v", res.StatusCode, res.Header, err))
			return
		}
		logger.Output(2, fmt.Sprintf("[status_code]=%d, [header]=%v", res.StatusCode, res.Header))
	}
}

func LogRequest(logger *slog.Logger) BeforeRequestHook {
	return func(r *http.Request, b []byte) {
		if logger == nil {
			logger = slog.Default()
		}
		if len(b) <= 1024 {
			logger.InfoContext(r.Context(), "send request", slog.String("method", r.Method), slog.String("request_url", r.URL.String()),
				slog.String("body", string(b)))
			return
		}
		logger.InfoContext(r.Context(), "send request", slog.String("method", r.Method), slog.String("request_url", r.URL.String()))
	}
}

func LogResponse(logger *slog.Logger) AfterRequestHook {
	return func(res *http.Response, err error) {
		if logger == nil {
			logger = slog.Default()
		}

		if err != nil {
			logger.ErrorContext(res.Request.Context(), "get response", slog.String("error", err.Error()),
				slog.Int("status_code", res.StatusCode))
			return
		}
		logger.InfoContext(res.Request.Context(), "get response", slog.Int("status_code", res.StatusCode))
	}
}
