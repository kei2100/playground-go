package strings

import "testing"

func TestConstantTimeCompare(t *testing.T) {
	if !ConstantTimeCompare("aaa", "aaa") {
		t.Error("unexpected false")
	}
	if ConstantTimeCompare("aaa", "bbb") {
		t.Error("unexpected true")
	}
}
