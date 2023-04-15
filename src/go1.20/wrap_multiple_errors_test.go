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

	wrap3 := errors.Join(nil, nil)
	assert.NoError(t, wrap3)

	wrap4 := errors.Join(err1, nil)
	assert.Error(t, wrap4)
	assert.True(t, errors.Is(wrap4, err1))

	wrap5 := errors.Join(nil, err1)
	assert.Error(t, wrap5)
	assert.True(t, errors.Is(wrap5, err1))
}
