package go1_25

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	var sum atomic.Int32
	var nums = []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, n := range nums {
		wg.Go(func() {
			sum.Add(n)
		})
	}

	wg.Wait()
	assert.EqualValues(t, 55, sum.Load())
}
