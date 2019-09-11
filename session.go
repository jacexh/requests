package requests

import (
	"crypto/tls"
	"encoding/json"
	"errors"
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

type (
	Option struct {
		Timeout            time.Duration
		AllowRedirects     bool
		InsecureSkipVerify bool
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
	return &Session{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: op.InsecureSkipVerify,
				},
			},
			Jar:     new(cookiejar.Jar),
			Timeout: op.Timeout,
		},
		op: op,
	}
}

func UnmarshalJSONResponse(v interface{}) Interceptor {
	return func(request *http.Request, response *http.Response, bytes []byte) error {
		if strings.HasPrefix(response.Header.Get("Content-Type"), "application/json") {
			return json.Unmarshal(bytes, v)
		}
		return errors.New("invalid content type")
	}
}

func (s *Session) Request(method string, path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	var err error
	var contentType string
	buff := GetBuffer()
	defer PutBuffer(buff)

	switch {
	case params.Body != nil:
		_, err = buff.Write(params.Body)

	case params.Json != nil:
		data, err := json.Marshal(params.Json)
		if err != nil {
			return nil, err
		}
		contentType = "application/json"
		if _, err = buff.Write(data); err != nil {
			return nil, err
		}

	case params.Files != nil:
		writer := multipart.NewWriter(buff)
		for k, v := range params.Data {
			if err := writer.WriteField(k, v); err != nil {
				return nil, err
			}
		}

		for field, fp := range params.Files {
			file, err := os.Open(fp)
			if err != nil {
				return nil, err
			}
			part, err := writer.CreateFormFile(field, filepath.Base(file.Name()))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		contentType = writer.FormDataContentType()

	case params.Data != nil:
		values := url.Values{}
		for k, v := range params.Data {
			values.Set(k, v)
		}
		buff.WriteString(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req, err := http.NewRequest(strings.ToUpper(method), path, buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "jacexh/requests")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	if params.Query != nil {
		query := req.URL.Query()
		for k, v := range params.Query {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	if interceptor != nil {
		err = interceptor(req, res, data)
	}
	return res, err
}

func (s *Session) Get(path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	return s.Request(http.MethodGet, path, params, op, interceptor)
}

func (s *Session) Post(path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	return s.Request(http.MethodPost, path, params, op, interceptor)
}

func (s *Session) Put(path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	return s.Request(http.MethodPut, path, params, op, interceptor)
}

func (s *Session) Patch(path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	return s.Request(http.MethodPatch, path, params, op, interceptor)
}

func (s *Session) Delete(path string, params Parameters, op Option, interceptor Interceptor) (*http.Response, error) {
	return s.Request(http.MethodDelete, path, params, op, interceptor)
}
