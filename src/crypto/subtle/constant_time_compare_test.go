package subtle

import (
	"crypto/subtle"
	"testing"
)

func TestConstantTimeCompare(t *testing.T) {
	a, b, c := "hello", "hello", "HELLO"
	const same = 1
	if subtle.ConstantTimeCompare([]byte(a), []byte(b)) != same {
		t.Error("unexpected")
	}
	if subtle.ConstantTimeCompare([]byte(a), []byte(c)) == same {
		t.Error("unexpected")
	}
}
