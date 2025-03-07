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
	Any map[string]any

	Params struct {
		Query   Any
		Data    Any
		Json    interface{}
		Files   Any
		Body    []byte
		Headers Any
	}

	Session struct {
		client        *http.Client
		userAgent     string
		globalHeaders Any
		beforeHooks   []BeforeRequestHook
		afterHooks    []AfterRequestHook
	}

	// BeforeRequestHook 会在调用`Client.Do(*http.Request)`前调用，此时已经完成了parama的自动装填
	BeforeRequestHook func(*http.Request, []byte)

	// AfterRequestHook 会在完成`Client.Do(*http.Request)`后立即调用
	AfterRequestHook func(*http.Response, error)

	// ResponseRender 用于对response的反序列化
	Unmarshaller func([]byte) error
)

func NewSession(opts ...Option) *Session {
	jar, _ := cookiejar.New(nil)
	s := &Session{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 300 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          500,
				MaxIdleConnsPerHost:   100,
				MaxConnsPerHost:       1000,
				IdleConnTimeout:       300 * time.Second,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 0,
			},
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		userAgent:     defaultUA,
		globalHeaders: make(Any),
		beforeHooks:   make([]BeforeRequestHook, 0),
		afterHooks:    make([]AfterRequestHook, 0),
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
			str, err := ToString(v)
			if err != nil {
				return "", err
			}
			if err := writer.WriteField(k, str); err != nil {
				break body
			}
		}

		for field, fp := range params.Files {
			str, err := ToString(fp)
			if err != nil {
				return "", err
			}
			file, err := os.Open(str)
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
			sv, err := ToString(v)
			if err != nil {
				return "", err
			}
			values.Set(k, sv)
		}
		_, err = w.Write([]byte(values.Encode()))
		contentType = "application/x-www-form-urlencoded"
	}

	return contentType, err
}

func (s *Session) Request(method, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	buff := GetBuffer()
	defer PutBuffer(buff)

	req, err := s.Prepare(context.Background(), method, path, params, buff)
	if err != nil {
		return nil, nil, err
	}
	return s.Send(req, unmarshaller)
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
		sv, err := ToString(v)
		if err != nil {
			return nil, err
		}
		req.Header.Set(k, sv)
	}
	req.Header.Set("User-Agent", s.userAgent)
	if req.Header.Get("Content-Type") == "" && autoContentType != "" {
		req.Header.Set("Content-Type", autoContentType)
	}
	if params.Headers != nil {
		for k, v := range params.Headers {
			sv, err := ToString(v)
			if err != nil {
				return nil, err
			}
			req.Header.Set(k, sv)
		}
	}

	if params.Query != nil {
		query := req.URL.Query()
		for k, v := range params.Query {
			sv, err := ToString(v)
			if err != nil {
				return nil, err
			}
			query.Set(k, sv)
		}
		req.URL.RawQuery = query.Encode()
	}

	for _, hook := range s.beforeHooks {
		hook(req, buff.Bytes())
	}
	return req, err
}

func (s *Session) Send(req *http.Request, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	var err error
	var data []byte

	res, err := s.client.Do(req)
	for _, hook := range s.afterHooks {
		hook(res, err)
	}
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return res, nil, err
	}
	_ = res.Body.Close()
	data = buf.Bytes()

	if unmarshaller != nil {
		err = unmarshaller(data)
	}
	return res, data, err
}

func (s *Session) RequestWithContext(ctx context.Context, method, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	buff := GetBuffer()
	defer PutBuffer(buff)

	req, err := s.Prepare(ctx, method, path, params, buff)
	if err != nil {
		return nil, nil, err
	}
	res, data, err := s.Send(req, unmarshaller)
	return res, data, err
}

func (s *Session) Get(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodGet, path, params, unmarshaller)
}

func (s *Session) GetWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodGet, path, params, unmarshaller)
}

func (s *Session) Post(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodPost, path, params, unmarshaller)
}

func (s *Session) PostWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPost, path, params, unmarshaller)
}

func (s *Session) Put(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodPut, path, params, unmarshaller)
}

func (s *Session) PutWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPut, path, params, unmarshaller)
}

func (s *Session) Patch(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodPatch, path, params, unmarshaller)
}

func (s *Session) PatchWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodPatch, path, params, unmarshaller)
}

func (s *Session) Delete(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodDelete, path, params, unmarshaller)
}

func (s *Session) DeleteWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodDelete, path, params, unmarshaller)
}

func (s *Session) Head(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodHead, path, params, unmarshaller)
}

func (s *Session) HeadWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodHead, path, params, unmarshaller)
}

func (s *Session) Connect(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodConnect, path, params, unmarshaller)
}

func (s *Session) ConnectWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodConnect, path, params, unmarshaller)
}

func (s *Session) Options(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodOptions, path, params, unmarshaller)
}

func (s *Session) OptionWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodOptions, path, params, unmarshaller)
}

func (s *Session) Trace(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.Request(http.MethodTrace, path, params, unmarshaller)
}

func (s *Session) TraceWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return s.RequestWithContext(ctx, http.MethodTrace, path, params, unmarshaller)
}
