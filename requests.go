package requests

import (
	"context"
	"net/http"
)

func Request(method, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Request(method, path, params, unmarshaller)
}

func RequestWithContext(ctx context.Context, method, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().RequestWithContext(ctx, method, path, params, unmarshaller)
}

func Get(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Get(path, params, unmarshaller)
}

func GetWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().GetWithContext(ctx, path, params, unmarshaller)
}

func Post(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Post(path, params, unmarshaller)
}

func PostWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().PostWithContext(ctx, path, params, unmarshaller)
}

func Put(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Put(path, params, unmarshaller)
}

func PutWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().PutWithContext(ctx, path, params, unmarshaller)
}

func Patch(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Patch(path, params, unmarshaller)
}

func PatchWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().PatchWithContext(ctx, path, params, unmarshaller)
}

func Delete(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Delete(path, params, unmarshaller)
}

func DeleteWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().DeleteWithContext(ctx, path, params, unmarshaller)
}

func Head(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Head(path, params, unmarshaller)
}

func HeadWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().HeadWithContext(ctx, path, params, unmarshaller)
}

func Connect(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Connect(path, params, unmarshaller)
}

func Options(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Options(path, params, unmarshaller)
}

func OptionWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().OptionWithContext(ctx, path, params, unmarshaller)
}

func Trace(path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().Put(path, params, unmarshaller)
}

func TraceWithContext(ctx context.Context, path string, params Params, unmarshaller Unmarshaller) (*http.Response, []byte, error) {
	return NewSession().TraceWithContext(ctx, path, params, unmarshaller)
}
