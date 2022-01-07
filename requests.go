package requests

import (
	"context"
	"net/http"
)

func Request(method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Request(method, path, params, interceptor)
}

func RequestWithContext(ctx context.Context, method, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().RequestWithContext(ctx, method, path, params, interceptor)
}

func Get(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Get(path, params, interceptor)
}

func GetWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().GetWithContext(ctx, path, params, interceptor)
}

func Post(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Post(path, params, interceptor)
}

func PostWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().PostWithContext(ctx, path, params, interceptor)
}

func Put(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Put(path, params, interceptor)
}

func PutWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().PutWithContext(ctx, path, params, interceptor)
}

func Patch(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Patch(path, params, interceptor)
}

func PatchWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().PatchWithContext(ctx, path, params, interceptor)
}

func Delete(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Delete(path, params, interceptor)
}

func DeleteWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().DeleteWithContext(ctx, path, params, interceptor)
}

func Head(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Head(path, params, interceptor)
}

func HeadWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().HeadWithContext(ctx, path, params, interceptor)
}

func Connect(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Connect(path, params, interceptor)
}

func Options(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Options(path, params, interceptor)
}

func OptionWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().OptionWithContext(ctx, path, params, interceptor)
}

func Trace(path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().Put(path, params, interceptor)
}

func TraceWithContext(ctx context.Context, path string, params Params, interceptor Interceptor) (*http.Response, []byte, error) {
	return NewSession().TraceWithContext(ctx, path, params, interceptor)
}
