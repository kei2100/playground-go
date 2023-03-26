package lang

import (
	"fmt"
	"testing"
)

type myCloser struct {
	message string
}

func (c *myCloser) Close() {
	fmt.Println(c.message)
}

func TestDeferCaptureVariable(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		c := &myCloser{message: "a"}
		defer c.Close()
		c = &myCloser{message: "aaa"}
		// Prints `a`
		// レシーバーや引数は defer の評価時にキャプチャされる（レシーバーも結局は引数）
	})
	t.Run("2", func(t *testing.T) {
		c := &myCloser{message: "b"}
		defer func() {
			c.Close()
		}()
		c = &myCloser{message: "bbb"}
		// Prints `bbb`
		// c は無名関数の引数になっていないのでキャプチャされない
	})
}
