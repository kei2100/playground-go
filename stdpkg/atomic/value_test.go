package atomic

import (
	"sync/atomic"
	"testing"
)

func TestValue(t *testing.T) {
	var v atomic.Value

	g := v.Load()
	if g != nil {
		t.Errorf("got %v, want nil", g)
	}

	v.Store("test")
	g = v.Load()
	if g.(string) != "test" {
		t.Errorf("got %v, want nil", g)
	}
}
