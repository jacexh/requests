package requests

import (
	"fmt"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestMergeOption(t *testing.T) {
	src := RequestOption{}
	src.Timeout = 3 * time.Second
	src.AllowRedirects = true

	target := RequestOption{}
	mergeOption(src, target)
	fmt.Printf("%v\n", target)
	fmt.Printf("%v\n", mergeOption(src, target))
}

func TestRequest(t *testing.T) {
	session := &Session{
		client: &fasthttp.Client{},
	}
	res := fasthttp.AcquireResponse()
	err := session.Request(
		"post", "https://en0393ftzur2k7.x.pipedream.net",
		RequestParameters{
			Query: map[string]string{"version": "2.0"},
			Data:  map[string]string{"hello": "world"},
			Files: map[string]string{"file": "README.md"},
		},
		RequestOption{},
		res,
		nil,
	)
	if err != nil {
		panic(err)
		t.Fail()
	}

	fmt.Println(string(res.Body()))
}
