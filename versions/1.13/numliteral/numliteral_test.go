package numliteral

import (
	"fmt"
)

func Example() {
	binary := 0b10
	fmt.Println(binary)
	fmt.Println(binary == 2)

	octal := 0o10
	fmt.Println(octal)
	fmt.Println(octal == 8)

	separated := 1_000
	fmt.Println(separated)
	fmt.Println(separated == 1000)

	// Output:
	// 2
	// true
	// 8
	// true
	// 1000
	// true
}
