package requests

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	Option struct {
		Name               string
		Timeout            time.Duration
		InsecureSkipVerify bool
		Headers            map[string]string
	}

	Parameters struct {
		Query  map[string]string
		Data   map[string]string
		Json   interface{}
		Files  map[string]string
		Body   []byte
		Header map[string]string
	}

	Session struct {
		client *http.Client
		op     Option
	}

	Interceptor func(*http.Request, *http.Response, []byte) error
)

func NewSession(op Option) *Session {
	jar, _ := cookiejar.New(nil)
	return &Session{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: op.InsecureSkipVerify,
				},
			},
			Jar:     jar,
			Timeout: op.Timeout,
		},
		op: op,
	}
}

func UnmarshalJSONResponse(v interface{}) Interceptor {
	return func(request *http.Request, response *http.Response, bytes []byte) error {
		return json.Unmarshal(bytes, v)
	}
}

func (s *Session) WithClient(client *http.Client) *Session {
	s.client = client
	s.client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = s.op.InsecureSkipVerify
	s.client.Timeout = s.op.Timeout
	return s
}

func (s *Session) WithOption(op Option) *Session {
	s.op = op
	if s.client != nil {
		s.client.Timeout = s.op.Timeout
		s.client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = s.op.InsecureSkipVerify
	}
	return s
}

func (s *Session) writeBody(params Parameters, w io.Writer) (string, error) {
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

func (s *Session) Prepare(method, path string, params Parameters, body io.ReadWriter) (*http.Request, error) {
	var err error
	var contentType string

	contentType, err = s.writeBody(params, body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(strings.ToUpper(method), path, body)
	if err != nil {
		return nil, err
	}

	// begin to set headers
	if s.op.Headers != nil {
		for k, v := range s.op.Headers {
			req.Header.Set(k, v)
		}
	}
	if s.op.Name == "" {
		req.Header.Set("User-Agent", defaultUA)
	} else {
		req.Header.Set("User-Agent", s.op.Name)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if params.Header != nil {
		for k, v := range params.Header {
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
	return req, err
}

func (s *Session) Send(req *http.Request, interceptor Interceptor) (*http.Response, []byte, error) {
	var err error
	var data []byte
	res, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	data, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return res, nil, err
	}
	_ = res.Body.Close()
	if interceptor != nil {
		err = interceptor(req, res, data)
	}
	return res, data, err
}

func (s *Session) Request(method, path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	buff := GetBuffer()
	defer PutBuffer(buff)

	req, err := s.Prepare(method, path, params, buff)
	if err != nil {
		return nil, nil, err
	}

	res, data, err := s.Send(req, interceptor)
	return res, data, err
}

func (s *Session) Get(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodGet, path, params, interceptor)
}

func (s *Session) Post(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPost, path, params, interceptor)
}

func (s *Session) Put(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPut, path, params, interceptor)
}

func (s *Session) Patch(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodPatch, path, params, interceptor)
}

func (s *Session) Delete(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodDelete, path, params, interceptor)
}

func (s *Session) Head(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodHead, path, params, interceptor)
}

func (s *Session) Connect(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodConnect, path, params, interceptor)
}

func (s *Session) Options(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodOptions, path, params, interceptor)
}

func (s *Session) Trace(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return s.Request(http.MethodTrace, path, params, interceptor)
}
