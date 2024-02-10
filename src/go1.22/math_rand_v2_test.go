package go1_22

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMathRandV2RandN(t *testing.T) {
	n := 100
	got := rand.N(n)
	assert.GreaterOrEqual(t, got, 0)
	assert.Less(t, got, n)
}
