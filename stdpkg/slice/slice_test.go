package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func _oooooExampleCap() {
	s1 := make([]struct{}, 0)
	s2 := append(s1, struct{}{})
	s3 := make([]struct{}, 0, 10)

	// capで内部の配列のサイズを検証する

	// Output:
	// 0
	// 1
	// 10
	fmt.Println(cap(s1))
	fmt.Println(cap(s2))
	fmt.Println(cap(s3))
}

func TestAppendNilSlice(t *testing.T) {
	var s1 []string
	s2 := make([]string, 0)

	s1 = append(s1, "test")
	s2 = append(s2, "test")

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("not same s1: %v, s2: %v", s1, s2)
	}
}
