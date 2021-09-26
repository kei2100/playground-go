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

func setupMulti(t *testing.T, ups ...func(*testing.T) func(*testing.T)) (teardownMulti func(*testing.T)) {
	t.Helper()

	tds := make([]func(*testing.T), len(ups))
	for i, up := range ups {
		tds[i] = up(t)
	}
	return func(t *testing.T) {
		t.Helper()
		for _, td := range tds {
			td(t)
		}
	}
}

func setupA(t *testing.T) func(*testing.T) {
	t.Helper()
	t.Log("setup A")
	return func(t *testing.T) {
		t.Helper()
		t.Log("teardown A")
	}
}

func setupB(t *testing.T) func(*testing.T) {
	t.Helper()
	t.Log("setup B")
	return func(t *testing.T) {
		t.Helper()
		t.Log("teardown B")
	}
}

func TestMulti(t *testing.T) {
	teardown := setupMulti(t, setupA, setupB)
	defer teardown(t)

	t.Log("do test")
}
