package maptest

import "testing"

func TestNilMap(t *testing.T) {
	var m map[int]int
	// m[0] = 0 // panics

	m = make(map[int]int)
	m[0] = 0
}
