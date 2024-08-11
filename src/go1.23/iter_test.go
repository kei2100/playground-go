package go1_23

import (
	"fmt"
	"iter"
	"testing"
)

// iter package に以下が定義されている
//
// * type Seq[V any] func(yield func(V) bool)
// * type Seq2[K, V any] func(yield func(K, V) bool)
// * func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {...}
// * func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {...}

// Seq は 「range over func」 に渡すことのできる `func(func(K) bool)` を表している
// Seq2 は 「range over func」 に渡すことのできる `func(func(K, V) bool)` を表している

func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				break
			}
		}
	}
}

func squares(n int) iter.Seq2[int, int64] {
	return func(yield func(int, int64) bool) {
		for i := range n {
			if !yield(i, int64(i)*int64(i)) {
				break
			}
		}
	}
}

func TestSeq(t *testing.T) {
	for i := range count(5) {
		t.Logf("i: %d", i)
		// i: 0
		// i: 1
		// i: 2
		// i: 3
		// i: 4
	}
	for i, i2 := range squares(5) {
		t.Logf("i: %d, i2: %d", i, i2)
		// i: 0, i2: 0
		// i: 1, i2: 1
		// i: 2, i2: 4
		// i: 3, i2: 9
		// i: 4, i2: 16
	}
}

// Pull, Pull2 は「プッシュ型」のイテレーターシーケンスを next() と stop() でアクセスする「プル型」のイテレーターに変換する

func TestPull(t *testing.T) {
	rangeFunc := iter.Seq[string](func(yield func(string) bool) {
		for _, s := range []string{"foo", "bar", "baz"} {
			if !yield(s) {
				break
			}
		}
	})

	next, stop := iter.Pull(rangeFunc)

	s, ok := next()
	fmt.Println(s, ok) // foo true

	s, ok = next()
	fmt.Println(s, ok) // bar true

	stop()

	s, ok = next()
	fmt.Println(s, ok) // <empty> false
}
