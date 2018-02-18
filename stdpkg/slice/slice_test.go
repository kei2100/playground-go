package slice

import "fmt"

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
