package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultUA = "jacexh/requests - a go client for human"
)

type (
	// Any 可用于query/headers/data/files传参
	Any map[string]string

	Params struct {
		Query   Any
		Data    Any
		Json    interface{}
		Files   Any
		Body    []byte
		Headers Any
	}

	Session struct {
		client          *http.Client
		userAgent       string
		globalHeaders   Any
		requestPrinter  RequestPrinter
		responsePrinter ResponsePrinter
	}

	Interceptor func(*http.Request, *http.Response, []byte) error
)

func UnmarshalJSONResponse(v interface{}) Interceptor {
	return func(request *http.Request, response *http.Response, bytes []byte) error {
		return json.Unmarshal(bytes, v)
	}
}

func NewSession(opts ...Option) *Session {
	jar, _ := cookiejar.New(nil)
	s := &Session{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   10,
				MaxConnsPerHost:       1000,
				IdleConnTimeout:       60 * time.Second,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		userAgent:     defaultUA,
		globalHeaders: make(Any),
	}

	s.Apply(opts...)
	return s
}

func (s *Session) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Session) renderBody(params Params, w io.Writer) (string, error) {
	var err error
	var contentType string

body:
	switch {
	case params.Body != nil:
		_, err = w.Write(params.Body)
		if err != nil {
			break body
		}

	case params.Json != nil:
		data, err := json.Marshal(params.Json)
		if err != nil {
			break body
		}
		contentType = "application/json"
		if _, err = w.Write(data); err != nil {
			break body
		}

	case params.Files != nil:
		writer := multipart.NewWriter(w)
		for k, v := range params.Data {
			if err := writer.WriteField(k, v); err != nil {
				break body
			}
		}

		for field, fp := range params.Files {
			file, err := os.Open(fp)
			if err != nil {
				break body
			}
			defer file.Close()

			part, err := writer.CreateFormFile(field, filepath.Base(file.Name()))
			if err != nil {
				break body
			}
			_, err = io.Copy(part, file)
			if err != nil {
				break body
			}
		}
		if err := writer.Close(); err != nil {
			break body
		}
		contentType = writer.FormDataContentType()

	case params.Data != nil:
		values := url.Values{}
		for k, v := range params.Data {
			values.Set(k, v)
		}
		_, err = w.Write([]byte(values.Encode()))
		contentType = "application/x-www-form-urlencoded"
	}

	return contentType, err
}

func (s *Session) Prepare(ctx context.Context, method, path string, params Params, buff *bytes.Buffer) (*http.Request, error) {
	var err error
	var autoContentType string
	var req *http.Request

	autoContentType, err = s.renderBody(params, buff)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequestWithContext(ctx, strings.ToUpper(method), path, buff)
	if err != nil {
		return nil, err
	}

	// begin to set headers
	for k, v := range s.globalHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", s.userAgent)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", autoContentType)
	}
	if params.Headers != nil {
		for k, v := range params.Headers {
			req.Header.Set(k, v)
		}
	}

	if params.Query != nil {
		query := req.URL.Query()
		for k, v := range params.Query {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	if s.requestPrinter != nil {
		s.requestPrinter.LogRequest(req, buff.Bytes())
	}
	return req, err
}

func (s *Session) Send(req *http.Request, interceptor Interceptor) (*http.Response, []byte, error) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	var err error
	var data []byte

	res, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return res, nil, err
	}
	_ = res.Body.Close()
	data = buf.Bytes()

	if s.responsePrinter != nil {
		s.responsePrinter.LogResponse(res, data)
	}

	if interceptor != nil {
		err = interceptor(req, res, data)
	}
	return res, data, err
}

func (s *Session) Request(method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	buff := GetBuffer()
	defer PutBuffer(buff)

	req, err := s.Prepare(context.Background(), method, path, params, buff)
	if err != nil {
		return nil, nil, err
	}
	res, data, err := s.Send(req, interceptor)
	return res, data, err
}

func (s *Session) RequestWithContext(ctx context.Context, method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	buff := GetBuffer()
	defer PutBuffer(buff)

	req, err := s.Prepare(ctx, method, path, params, buff)
	if err != nil {
		return nil, nil, err
	}
	res, data, err := s.Send(req, interceptor)
	return res, data, err
}

func (s *Session) Get(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodGet, path, params, interceptor)
}

func (s *Session) GetWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodGet, path, params, interceptor)
}

func (s *Session) Post(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPost, path, params, interceptor)
}

func (s *Session) PostWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPost, path, params, interceptor)
}

func (s *Session) Put(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPut, path, params, interceptor)
}

func (s *Session) PutWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPut, path, params, interceptor)
}

func (s *Session) Patch(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPatch, path, params, interceptor)
}

func (s *Session) PatchWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPatch, path, params, interceptor)
}

func (s *Session) Delete(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodDelete, path, params, interceptor)
}

func (s *Session) DeleteWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodDelete, path, params, interceptor)
}

func (s *Session) Head(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodHead, path, params, interceptor)
}

func (s *Session) HeadWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodHead, path, params, interceptor)
}

func (s *Session) Connect(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodConnect, path, params, interceptor)
}

func (s *Session) ConnectWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodConnect, path, params, interceptor)
}

func (s *Session) Options(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodOptions, path, params, interceptor)
}

func (s *Session) OptionWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodOptions, path, params, interceptor)
}

func (s *Session) Trace(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodTrace, path, params, interceptor)
}

func (s *Session) TraceWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodTrace, path, params, interceptor)
}
