package requests

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestInterceptor(t *testing.T) {
	session := NewSession()
	ret := make(map[string]interface{})
	_, _, err := session.Request(
		"get", "https://httpbin.org/json",
		Params{},
		UnmarshalJSONResponse(&ret),
	)

	if err != nil {
		t.FailNow()
	}
	if len(ret) == 0 {
		t.FailNow()
	}

	if _, ok := ret["slideshow"]; !ok {
		t.FailNow()
	}
}

func TestCreateBinOnRequestBin(t *testing.T) {
	session := NewSession(
		WithUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"),
		WithGlobalTimeout(90*time.Second),
		WithStdRequestPrinter(),
		WithStdResponsePrinter(),
	)
	_, _, err := session.Get("https://requestbin.net", Params{}, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	obj := make(map[string]interface{})
	_, data, err := session.Post("https://requestbin.net/api/v1/bins", Params{
		Data:    Any{"private": "false"},
		Headers: Any{"Origin": "https://requestbin.net"},
	},
		UnmarshalJSONResponse(&obj))
	if err != nil {
		if data != nil {
			t.Fatalf("%s\n%s", err.Error(), data)
		} else {
			t.Fatal(err.Error())
		}
	}
	if _, ok := obj["name"]; !ok {
		t.Fatal(string(data))
	}
	bin := obj["name"].(string)
	res, _, err := session.Get("https://requestbin.net/r/"+bin, Params{Data: Any{"hello": "world"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.FailNow()
	}
	log.Printf("open %s for more details", "https://requestbin.net/r/"+bin+"?inspect")

	_, _, err = session.Post("https://requestbin.net/r/"+bin, Params{
		Query: Any{"format": "json"},
		Json:  map[string]interface{}{"hello": "foobar", "version": 1}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, _, err = session.Post("https://requestbin.net/r/"+bin, Params{
		Query: Any{"format": "multipart"},
		Data:  Any{"version": "3"},
		Files: Any{"file": "README.md"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, _, err = session.Post("https://requestbin.net/r/"+bin, Params{Body: []byte(`i am body`)}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestRequestWithContext(t *testing.T) {
	session := NewSession(
		WithUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"),
		WithGlobalTimeout(30*time.Second),
	)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	_, _, err := session.PostWithContext(ctx, "https://requestbin.net/ip", Params{Data: Any{"fizz": "buzz"}}, nil)
	if err == nil {
		t.Fatal("deadline did not exceeded")
	}
}
