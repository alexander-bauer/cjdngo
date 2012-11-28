package admin

import (
	"testing"
)

func TestPing(t *testing.T) {
	if !Ping("127.0.0.1", "11234") {
		t.Fail()
	}
}

func TestCookie(t *testing.T) {
	if len(Cookie("127.0.0.1", "11234")) == 0 {
		t.Fail()
	}
}
