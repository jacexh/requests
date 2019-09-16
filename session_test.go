package requests

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func TestInterceptor(t *testing.T) {
	session := NewSession(Option{})
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
}

func TestCreateBinOnRequestBin(t *testing.T) {
	session := NewSession(Option{
		Name:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36",
		Timeout: 90 * time.Second,
	})
	_, _, err := session.Get("http://requestbin.net", Params{}, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	obj := make(map[string]interface{})
	_, data, err := session.Post("http://requestbin.net/api/v1/bins", Params{
		Data:   Any{"private": "false"},
		Header: Any{"Origin": "http://requestbin.net"},
	},
		UnmarshalJSONResponse(&obj))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if _, ok := obj["name"]; !ok {
		t.Fatal(string(data))
	}
	bin := obj["name"].(string)
	res, _, err := session.Get("http://requestbin.net/r/"+bin, Params{Data: Any{"hello": "world"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.FailNow()
	}
	log.Printf("open %s for more details", "http://requestbin.net/r/"+bin+"?inspect")
}
