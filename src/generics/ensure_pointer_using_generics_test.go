package generics

import "testing"

func TestJSONUnmarshal(t *testing.T) {
	type something struct {
		Foo string `json:"foo"`
	}
	json := []byte(`{"foo": "bar"}`)
	var dst something

	if err := JSONUnmarshal(json, &dst); err != nil {
		t.Error(err)
	}
	if g, w := dst.Foo, "bar"; g != w {
		t.Errorf("\ngot :%v\nwant:%v", dst.Foo, "bar")
	}
	// JSONUnmarshal(json, dst)
	// のような呼び出しはコンパイルエラーとなり、呼び出し側にポインタを使うことを強制できる。
}
