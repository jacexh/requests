package requests

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Option func(*Session)

func WithClient(client *http.Client) Option {
	return func(s *Session) {
		if client == nil {
			return
		}
		s.client = client
		if s.client.Jar == nil {
			jar, _ := cookiejar.New(nil)
			s.client.Jar = jar
		}
	}
}

func WithGlobalTimeout(t time.Duration) Option {
	return func(s *Session) {
		if s.client == nil {
			return
		}
		s.client.Timeout = t
	}
}

func WithUserAgent(name string) Option {
	return func(s *Session) {
		s.userAgent = name
	}
}

func WithGlobalHeader(header Any) Option {
	return func(s *Session) {
		s.globalHeaders = make(Any)
		for k, v := range header {
			s.globalHeaders[k] = v
		}
	}
}

func WithTransport(t http.RoundTripper) Option {
	return func(s *Session) {
		if t == nil {
			return
		}
		s.client.Transport = t
	}
}

func WithBeforeHooks(hooks ...BeforeRequestHook) Option {
	return func(s *Session) {
		s.beforeHooks = append(s.beforeHooks, hooks...)
	}
}

func WithAfterHooks(hooks ...AfterRequestHook) Option {
	return func(s *Session) {
		s.afterHooks = append(s.afterHooks, hooks...)
	}
}
