package strings

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestHeadBytes(t *testing.T) {
	tt := []struct {
		s    string
		size int
		want string
	}{
		{
			s:    "あい",
			size: 2,
			want: "",
		},
		{
			s:    "あい",
			size: 3,
			want: "あ",
		},
		{
			s:    "あい",
			size: 4,
			want: "あ",
		},
		{
			s:    "あい",
			size: 5,
			want: "あ",
		},
		{
			s:    "あい",
			size: 6,
			want: "あい",
		},
		{
			s:    "あい",
			size: 7,
			want: "あい",
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := headBytes(te.s, te.size)
			if g, w := got, te.want; g != w {
				t.Errorf("got %v, want %v", g, w)
			}
		})
	}
}

func headBytes(s string, size int) string {
	// バイトサイズ以内になるように切り取り
	if len(s) <= size {
		return s
	}
	var i int
	runes := []rune(s)
	for _, r := range runes {
		size -= utf8.RuneLen(r)
		if size < 0 {
			break
		}
		i++
	}
	return string(runes[:i])
}
