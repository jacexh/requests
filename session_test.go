package requests_test

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacexh/requests"
)

type Response struct {
	Headers map[string]string `json:"headers,omitempty"`
	Queries map[string]string `json:"queries,omitempty"`
	Body    string            `json:"body,omitempty"`
	Error   string            `json:"error,omitempty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		Headers: make(map[string]string),
		Queries: make(map[string]string),
	}
	// copy request headers
	for k := range r.Header {
		res.Headers[k] = r.Header.Get(k)
	}
	// copy query params
	for k := range r.URL.Query() {
		res.Queries[k] = r.URL.Query().Get(k)
	}
	// copy request body
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		res.Error = err.Error()
	} else if len(data) > 0 {
		res.Body = base64.StdEncoding.EncodeToString(data)
	}

	data, _ = json.Marshal(res)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func TestInterceptor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	session := requests.NewSession()
	session.Apply(
		requests.WithBeforeHooks(requests.LogRequest(slog.Default())),
		requests.WithAfterHooks(requests.LogResponse(slog.Default())),
	)
	ret := new(Response)
	payload := []byte("hello world")
	_, _, err := session.Request(
		"Post", ts.URL,
		requests.Params{Body: payload},
		requests.UnmarshalJSON(&ret),
	)

	if err != nil {
		t.FailNow()
	}
	if ret.Headers["User-Agent"] != "jacexh/requests - a go client for human" {
		t.FailNow()
	}

	slog.Info("print response", slog.Any("response", ret))

	expectedBody := base64.StdEncoding.EncodeToString(payload)
	if expectedBody != ret.Body {
		t.FailNow()
	}
}
