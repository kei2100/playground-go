package argon2

import (
	"fmt"
	"testing"
)

func TestGenerateFromPassword(t *testing.T) {
	g, _ := GenerateFromPassword("foo")
	fmt.Println(g)
	err := CompareHashAndPassword(g, "foo")
	fmt.Println(err)
	err = CompareHashAndPassword(g, "bar")
	fmt.Println(err)
}
