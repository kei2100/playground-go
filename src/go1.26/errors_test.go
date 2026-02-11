package go1_26

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myError struct {
	msg string
}

func (e *myError) Error() string {
	return e.msg
}

func oops() error {
	return &myError{msg: "this is my error"}
}

func TestAsError(t *testing.T) {
	err := oops()
	merr, ok := errors.AsType[*myError](err)
	assert.True(t, ok)
	fmt.Println(merr.msg)
}
