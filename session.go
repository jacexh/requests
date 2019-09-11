package requests

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	RequestOption struct {
		Timeout        time.Duration
		AllowRedirects bool
	}

	RequestParameters struct {
		Query  map[string]string
		Data   map[string]string
		Json   interface{}
		Files  map[string]string
		Body   []byte
		Header map[string]string
	}

	Session struct {
		client *fasthttp.Client
		op     *RequestParameters
		cookie *fasthttp.Cookie
	}
)

func mergeOption(src, target RequestOption) RequestOption {
	if target.Timeout == 0 {
		target.Timeout = src.Timeout
	}
	if target.AllowRedirects == false {
		target.AllowRedirects = src.AllowRedirects
	}
	return target
}

func NewSession() *Session {
	return &Session{client: &fasthttp.Client{}}
}

func (s *Session) Request(method string, path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var err error
	req.Header.SetMethod(strings.ToUpper(method))
	req.URI().Update(path)

	if params.Query != nil {
		val := url.Values{}
		for k, v := range params.Query {
			val.Set(k, v)
		}
		req.URI().SetQueryString(val.Encode())
	}

	switch {
	case params.Body != nil:
		req.SetBody(params.Body)

	case params.Json != nil:
		data, err := json.Marshal(params.Json)
		if err != nil {
			return err
		}
		req.Header.SetContentType("application/json; charset=UTF-8")
		req.SetBody(data)

	case params.Files != nil:
		writer := multipart.NewWriter(req.BodyWriter())
		for k, v := range params.Data {
			if err := writer.WriteField(k, v); err != nil {
				return err
			}
		}

		for field, fp := range params.Files {
			file, err := os.Open(fp)
			if err != nil {
				return err
			}
			part, err := writer.CreateFormFile(field, filepath.Base(file.Name()))
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
		}
		if err := writer.Close(); err != nil {
			return err
		}
		req.Header.SetContentType(writer.FormDataContentType())

	case params.Data != nil:
		values := url.Values{}
		for k, v := range params.Data {
			values.Set(k, v)
		}
		req.SetBodyString(values.Encode())
		req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	}

	err = s.client.Do(req, res)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Get(path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	return s.Request(fasthttp.MethodGet, path, params, op, res, interceptor)
}

func (s *Session) Post(path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	return s.Request(fasthttp.MethodPost, path, params, op, res, interceptor)
}

func (s *Session) Put(path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	return s.Request(fasthttp.MethodPut, path, params, op, res, interceptor)
}

func (s *Session) Patch(path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	return s.Request(fasthttp.MethodPatch, path, params, op, res, interceptor)
}

func (s *Session) Delete(path string, params RequestParameters, op RequestOption, res *fasthttp.Response, interceptor interface{}) error {
	return s.Request(fasthttp.MethodDelete, path, params, op, res, interceptor)
}
