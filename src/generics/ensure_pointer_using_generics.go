package generics

import (
	"encoding/json"
	"fmt"
)

// JSONUnmarshal は、型パラメータ T を定義し引数を dst *T とすることで、T のポインタを使用することを呼び出し側に強制しています。
// 従来はこのような引数は interface{} を使って定義されておりどのような型でも渡すことができていたため、
// ポインタでない型を渡してしまうことによる実行時エラーの懸念がありました。
func JSONUnmarshal[T any](data []byte, dst *T) error {
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("generics: json unmarshal: %w", err)
	}
	return nil
}
