package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	src := `package main
import "fmt"
func main() {
fmt.Println("Hello")
}
`
	want := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`

	got, err := Format([]byte(src))
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}
