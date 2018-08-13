package errors

import (
	"errors"
	"testing"
)

func newe(s string) error {
	return errors.New(s)
}

func TestJoin(t *testing.T) {
	tt := []struct {
		errs []error
		sep  string
		want string
	}{
		{
			errs: nil,
			sep:  ":",
			want: "",
		},
		{
			errs: []error{newe("a")},
			sep:  ":",
			want: "a",
		},
		{
			errs: []error{newe("a"), newe("b")},
			sep:  ":",
			want: "a:b",
		},
		{
			errs: []error{newe("a"), newe("b"), newe("c")},
			sep:  ",",
			want: "a,b,c",
		},
	}
	for i := range tt {
		got := Join(tt[i].errs, tt[i].sep)
		if g, w := got, tt[i].want; g != w {
			t.Errorf("tt[%v] got %v, want %v", i, g, w)
		}
	}
}
