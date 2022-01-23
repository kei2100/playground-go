package generics

import (
	"testing"
)

func TestGenericsFunctions(t *testing.T) {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	t.Run("non-generics functions", func(t *testing.T) {
		if g, w := SumInts(ints), int64(46); g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumInts(ints), 46)
		}
		if g, w := SumFloats(floats), 62.97; g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumFloats(floats), 62.97)
		}
	})
	t.Run("generics function with type parameters", func(t *testing.T) {
		// use type parameters such as `[string, int64]`
		if g, w := SumIntsOrFloats[string, int64](ints), int64(46); g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumInts(ints), 46)
		}
		if g, w := SumIntsOrFloats[string, float64](floats), 62.97; g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumFloats(floats), 62.97)
		}
	})
	t.Run("generics function without type parameters", func(t *testing.T) {
		// 可能な場合、Go コンパイラは type parameter の推論を行うため、呼び出しコードでの type parameter の記述を省略することができる
		// ※ 引数のない generic 関数などは推論ができないので常に可能なわけではない
		if g, w := SumIntsOrFloats(ints), int64(46); g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumInts(ints), 46)
		}
		if g, w := SumIntsOrFloats(floats), 62.97; g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumFloats(floats), 62.97)
		}
	})
	t.Run("use declared type constraint", func(t *testing.T) {
		if g, w := SumNumbers(ints), int64(46); g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumInts(ints), 46)
		}
		if g, w := SumNumbers(floats), 62.97; g != w {
			t.Errorf("\ngot :%v\nwant:%v", SumFloats(floats), 62.97)
		}
	})
}
