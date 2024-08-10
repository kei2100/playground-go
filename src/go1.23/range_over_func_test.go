package go1_23

import (
	"fmt"
	"testing"
)

// 以下の関数を for range することができるようになった
// func(func() bool)
// func(func(K) bool)
// func(func(K, V) bool)

func TestRangeOverFunc1(t *testing.T) {
	// func(func() bool)
	rangeFunc := func(yield func() bool) {
		ret := yield()
		fmt.Printf("yield 1: ret: %v\n", ret)
		if !ret {
			return
		}

		ret = yield()
		fmt.Printf("yield 2: ret: %v\n", ret)
		if !ret {
			return
		}

		ret = yield()
		fmt.Printf("yield 3: ret: %v\n", ret)
		if !ret {
			return
		}
	}
	for range rangeFunc {
		fmt.Printf("for 1: ")
		// for 1: yield 1: ret: true
		// for 1: yield 2: ret: true
		// for 1: yield 3: ret: true
	}
	for range rangeFunc {
		fmt.Printf("for 2: ")
		break
		// for 2: yield 1: ret: false
	}
}

func TestRangeOverFunc2(t *testing.T) {
	// func(func(K) bool)
	rangeFunc := func(yield func(int) bool) {
		for i := range 5 {
			if !yield(i) {
				return
			}
		}
	}
	for i := range rangeFunc {
		fmt.Printf("for: %d\n", i)
		// for: 0
		// for: 1
		// for: 2
		// for: 3
		// for: 4
	}
}

func TestRangeOverFunc3(t *testing.T) {
	// func(func(K, V) bool)
	rangeFunc := func(yield func(int, string) bool) {
		for i := range 5 {
			if !yield(i, fmt.Sprintf("value %d", i)) {
				return
			}
		}
	}
	for i, v := range rangeFunc {
		fmt.Printf("for: %d %s\n", i, v)
		//for: 0 value 0
		//for: 1 value 1
		//for: 2 value 2
		//for: 3 value 3
		//for: 4 value 4
	}
}
