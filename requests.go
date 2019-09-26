package requests

import (
	"context"
	"net/http"
)

func Request(method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Request(method, path, params, interceptor)
}

func RequestWithContext(ctx context.Context, method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).RequestWithContext(ctx, method, path, params, interceptor)
}

func Get(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Get(path, params, interceptor)
}

func GetWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).GetWithContext(ctx, path, params, interceptor)
}

func Post(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Post(path, params, interceptor)
}

func PostWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).PostWithContext(ctx, path, params, interceptor)
}

func Put(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Put(path, params, interceptor)
}

func PutWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).PutWithContext(ctx, path, params, interceptor)
}

func Patch(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Patch(path, params, interceptor)
}

func PatchWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).PatchWithContext(ctx, path, params, interceptor)
}

func Delete(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Delete(path, params, interceptor)
}

func DeleteWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).DeleteWithContext(ctx, path, params, interceptor)
}

func Head(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Head(path, params, interceptor)
}

func HeadWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).HeadWithContext(ctx, path, params, interceptor)
}

func Connect(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Connect(path, params, interceptor)
}

func Options(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Options(path, params, interceptor)
}

func OptionWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).OptionWithContext(ctx, path, params, interceptor)
}

func Trace(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).Put(path, params, interceptor)
}

func TraceWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession(Option{}).TraceWithContext(ctx, path, params, interceptor)
}
