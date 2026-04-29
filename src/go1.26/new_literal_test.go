package go1_26

import (
	"fmt"
	"testing"
)

func TestNewLiteral(t *testing.T) {
	type S struct {
		a *int
		b *bool
	}
	s := &S{
		a: new(10),
		b: new(true),
	}
	fmt.Printf("%#v\n", s)
}
