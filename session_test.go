package requests

import (
	"fmt"
	"testing"
)

func TestRequest(t *testing.T) {
	session := NewSession(Option{})
	ret := make(map[string]interface{})
	res, err := session.Request(
		"get", "https://httpbin.org/cookies",
		Parameters{},
		Option{},
		UnmarshalJSONResponse(&ret),
	)
	if err != nil {
		panic(err)
		t.Fail()
	}

	fmt.Println(ret)
	for k := range res.Header {
		fmt.Println(k, res.Header.Get(k))
	}
}
