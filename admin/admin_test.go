package admin

import (
	"testing"
)

func TestPing(t *testing.T) {
	if !Ping("127.0.0.1", "11234") {
		t.Fail()
	}
}
