package abs

import (
	"math"
	"testing"
)

func abs32(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}

func abs64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func TestAbs32(t *testing.T) {
	if g, w := abs32(math.MaxInt32), int32(math.MaxInt32); g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
	if g, w := abs32(-math.MaxInt32), int32(math.MaxInt32); g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
	if g, w := abs32(math.MinInt32+1), int32(math.MaxInt32); g != w {
		t.Errorf(" got %v, want %v", g, w)
	}
}
