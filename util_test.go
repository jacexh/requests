package requests_test

import (
	"log"
	"log/slog"
	"testing"

	"github.com/jacexh/requests"
)

func TestToString(t *testing.T) {
	values := []any{true, uint8(250), uint16(250), uint32(250), uint64(250)}
	for _, value := range values {
		log.Println(value)
		str, err := requests.ToString(value)
		if err != nil {
			slog.Error(err.Error())
			t.FailNow()
		}
		slog.Info("message", slog.String("value", str))
	}
}
