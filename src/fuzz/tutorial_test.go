package main

import (
	"testing"
	"unicode/utf8"
)

func Reverse(s string) string {
	//b := []byte(s)
	//for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
	//	b[i], b[j] = b[j], b[i]
	//}
	//return string(b)
	return ReverseFixed(s)
}

func ReverseFixed(s string) string {
	//r := []rune(s)
	//for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
	//	r[i], r[j] = r[j], r[i]
	//}
	//return string(r)
	return ReverseFixed2(s)
}

func ReverseFixed2(s string) string {
	if !utf8.ValidString(s) {
		return s
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	for _, tc := range testcases {
		rev := Reverse(tc.in)
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

// go1.18beta1 test -fuzz=Fuzz -fuzztime 60s github.com/kei2100/playground-go/src/fuzz
// のようにすると新たに fuzz test を実行する。
//
// エラーを発見すると、そのテストデータが testdata/fuzz/FuzzReverse に出力される。
// テストデータがある場合は、-fuzz フラグなしでもそのテストデータを使ったテストが実行される。
func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
