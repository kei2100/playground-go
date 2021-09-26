package algo

import (
	"fmt"
	"testing"
)

func TestSearchReciprocal(t *testing.T) {
	a := uint32(0x1ca7bc5b)
	b := uint32(0x6b5f13d3)
	v := uint32(1)

	for h := 0; h < 10; {
		a++
		b--
		if v*a*b == 1 {
			h++
			fmt.Printf("a: 0x%x, b 0x%x\n", a, b)
		}
	}
}

func TestScramble(t *testing.T) {
	a := uint32(0xdca7bc5b)
	b := uint32(0xab5f13d3)
	m := make(map[uint32]struct{})

	for v := uint32(1); v < 1000000; v++ {
		r := scramble(v, a, b)
		if _, ok := m[r]; ok {
			t.Errorf("duplicate error. v: %v, r :%v, a: 0x%x, b: 0x%x", v, r, a, b)
		}
		m[r] = struct{}{}
	}
	fmt.Println("done")
}

func scramble(v, a, b uint32) uint32 {
	v *= a

	v = ((v >> 1) & 0x55555555) | ((v & 0x55555555) << 1)
	v = ((v >> 2) & 0x33333333) | ((v & 0x33333333) << 2)
	v = ((v >> 4) & 0x0F0F0F0F) | ((v & 0x0F0F0F0F) << 4)
	v = ((v >> 8) & 0x00FF00FF) | ((v & 0x00FF00FF) << 8)
	v = (v >> 16) | (v << 16)

	v *= b
	return v
}
