package requests

import (
	"testing"
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
