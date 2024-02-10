package go1_22

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForRangeInteger(t *testing.T) {
	var got []int
	for i := range 3 {
		got = append(got, i)
	}
	assert.Equal(t, []int{0, 1, 2}, got)
}
