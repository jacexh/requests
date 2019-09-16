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
		Parameters{},
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
	_, _, err := session.Get("http://requestbin.net", Parameters{}, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	obj := make(map[string]interface{})
	_, data, err := session.Post("http://requestbin.net/api/v1/bins", Parameters{
		Data:   map[string]string{"private": "false"},
		Header: map[string]string{"Origin": "http://requestbin.net"},
	},
		UnmarshalJSONResponse(&obj))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if _, ok := obj["name"]; !ok {
		t.Fatal(string(data))
	}
	bin := obj["name"].(string)
	res, _, err := session.Get("http://requestbin.net/r/"+bin, Parameters{Data: map[string]string{"hello": "world"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.FailNow()
	}
	log.Printf("open %s for more details", "http://requestbin.net/r/"+bin+"?inspect")
}

func TestRequest(t *testing.T) {
	session := NewSession(Option{})
	_, _, err := session.Post("http://requestbin.net/r/1dwh0311", Parameters{Data: map[string]string{"version": "2.0"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, _, err = session.Post("http://requestbin.net/r/1dwh0311", Parameters{Json: map[string]string{"version": "20"}}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
}
