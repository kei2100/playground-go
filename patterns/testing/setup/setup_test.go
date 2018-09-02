package setup

import "testing"

var var1 = "var1"

func setup(t *testing.T) (teardown func(*testing.T)) {
	t.Helper()
	bk := var1
	var1 = "mock1"

	return func(t *testing.T) {
		t.Helper()
		var1 = bk
	}
}

func TestVar(t *testing.T) {
	t.Run("use setup", func(t *testing.T) {
		teardown := setup(t)
		defer teardown(t)

		if g, w := var1, "mock1"; g != w {
			t.Errorf("var1 got %v, want %v", g, w)
		}
	})

	t.Run("not use setup", func(t *testing.T) {
		if g, w := var1, "var1"; g != w {
			t.Errorf("var1 got %v, want %v", g, w)
		}
	})
}
