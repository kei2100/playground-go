package setup

import "testing"

var var1 = "var1"
var var2 = "var2"

func WithMockVar1(t *testing.T, f func(*testing.T)) {
	bk := var1
	var1 = "mock1"
	defer func() { var1 = bk }()

	f(t)
}

func WithMockVar2(t *testing.T, f func(*testing.T)) {
	bk := var2
	var2 = "mock2"
	defer func() { var2 = bk }()

	f(t)
}

func TestVar(t *testing.T) {
	WithMockVar1(t, func(t *testing.T) {
		WithMockVar2(t, func(t *testing.T) {
			if g, w := var1, "mock1"; g != w {
				t.Errorf("var1 got %v, want %v", g, w)
			}

			if g, w := var2, "mock2"; g != w {
				t.Errorf("var2 got %v, want %v", g, w)
			}
		})
	})

	t.Error("to be improve")
}
