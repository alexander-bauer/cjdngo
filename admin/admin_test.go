package admin

import (
	"testing"
)

func TestPing(t *testing.T) {
	Ping("127.0.0.1", "11234")
}
