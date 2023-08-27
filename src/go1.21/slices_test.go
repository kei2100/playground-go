package go1_21

import (
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestSlicesBinarySearch(t *testing.T) {
	// slices.BinarySearch
	ss := []string{"a", "b", "c", "d", "e", "f"}
	i := slices.BinarySearch(ss, "d")
	if g, w := i, 3; g != w {
		t.Errorf("\ngot :%v\nwant:%v", i, 3)
	}
	// slices.BinarySearchFunc
	i = slices.BinarySearchFunc(ss, func(s string) bool {
		return strings.ToUpper(s) == "D"
	})
	if g, w := i, 3; g != w {
		t.Errorf("\ngot :%v\nwant:%v", i, 3)
	}
}

func TestSlicesContains_Index(t *testing.T) {
	// slices.Contains
	ss := []string{"a", "b", "c", "d", "e", "f"}
	got := slices.Contains(ss, "d")
	if g, w := got, true; g != w {
		t.Errorf("\ngot :%v\nwant:%v", got, true)
	}
	// slices.Index
	i := slices.Index(ss, "d")
	if g, w := i, 3; g != w {
		t.Errorf("\ngot :%v\nwant:%v", i, 3)
	}
	// slices.IndexFunc
	i = slices.IndexFunc(ss, func(s string) bool {
		return strings.ToUpper(s) == "D"
	})
	if g, w := i, 3; g != w {
		t.Errorf("\ngot :%v\nwant:%v", i, 3)
	}
}
