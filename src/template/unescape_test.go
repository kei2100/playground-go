package template

import "testing"

func TestUnescape(t *testing.T) {
	got := Unescape()
	want := "<html><body><p>Hello</p></body></html>"

	if g, w := got, want; g != w {
		t.Errorf("\ngot :%v\nwant:%v", got, want)
	}
}
