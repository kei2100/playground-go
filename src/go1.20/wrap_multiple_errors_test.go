package go1_20

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	wrap1 := errors.Join(err1, err2)
	wrap2 := fmt.Errorf("foo: bar: %w, %w", err1, err2)

	assert.True(t, errors.Is(wrap1, err1))
	assert.True(t, errors.Is(wrap1, err2))
	assert.True(t, errors.Is(wrap2, err1))
	assert.True(t, errors.Is(wrap2, err2))
}
