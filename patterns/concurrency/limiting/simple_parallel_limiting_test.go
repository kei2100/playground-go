package limiting

import (
	"testing"
	"fmt"
	"sync"
	"time"
)

func TestSimpleParallelLimiting(t *testing.T) {
	wg := new(sync.WaitGroup)
	limit := make(chan struct{}, 2)

	task := func(num int) {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		fmt.Printf("%v done\n", num)
		<- limit
	}

	const taskCount = 10
	wg.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		limit <- struct{}{}
		go task(i)
	}

	wg.Wait()
}
