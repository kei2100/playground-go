package go1_26

import "fmt"

//go:fix inline
func Old(s string) string {
	return New(s, "world")
}

func New(s string, ss string) string {
	return fmt.Sprintf("%s %s", s, ss)
}

func Foo() {
	// go fix　./... をすると Old 関数の呼び出し箇所が New 関数の呼び出しに置き換わる
	s := Old("hello")
	fmt.Println(s)
}
