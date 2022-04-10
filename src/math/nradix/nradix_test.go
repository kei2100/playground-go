package nradix

import (
	"fmt"
	"testing"
)

// go test -fuzz=Fuzz -fuzztime 10s github.com/kei2100/playground-go/src/math/nradix
func FuzzNRadix_ConvertToString_Binary(f *testing.F) {
	binary := New("01")
	hex := New("0123456789abcdef")
	corpus := []int64{2, 256, 0, -256}
	for _, c := range corpus {
		f.Add(c)
	}
	f.Fuzz(func(t *testing.T, in int64) {
		got := binary.ConvertToString(in)
		want := fmt.Sprintf("%b", in)
		if got != want {
			t.Errorf("binary: in: %d\ngot :%v\nwant:%v", in, got, want)
		}
		got = hex.ConvertToString(in)
		want = fmt.Sprintf("%x", in)
		if got != want {
			t.Errorf("hex: in: %d\ngot :%v\nwant:%v", in, got, want)
		}
	})
}
