package integer

import "testing"

func TestPtr(t *testing.T) {
	if g, w := *Ptr(1), 1; g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
}
