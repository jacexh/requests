package requests

import "net/http"

func Get(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Get(path, params, interceptor)
}

func Post(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Post(path, params, interceptor)
}

func Put(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Put(path, params, interceptor)
}

func Patch(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Patch(path, params, interceptor)
}

func Delete(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Delete(path, params, interceptor)
}

func Head(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Head(path, params, interceptor)
}

func Connect(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Connect(path, params, interceptor)
}

func Options(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Options(path, params, interceptor)
}

func Trace(path string, params Parameters, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Put(path, params, interceptor)
}
