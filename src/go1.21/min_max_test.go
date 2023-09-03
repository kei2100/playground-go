package go1_21

import "testing"

func TestMinMax(t *testing.T) {
	i := min(5, 4, 3, 1, 2)
	if g, w := i, 1; g != w {
		t.Errorf("\ngot :%v\nwant:%v", i, 1)
	}
	s := max("a", "あ", "A")
	if g, w := s, "あ"; g != w {
		t.Errorf("\ngot :%v\nwant:%v", s, "あ")
	}
}
