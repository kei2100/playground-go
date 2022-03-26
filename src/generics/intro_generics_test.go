package generics

import "testing"

func TestGMin(t *testing.T) {
	// Providing the type argument to GMin, in this case int, is called `instantiation`.
	instantiated := GMin[int]
	g := instantiated(11, 10)
	if g, w := g, 10; g != w {
		t.Errorf("\ngot :%v\nwant:%v", g, 10)
	}
}
