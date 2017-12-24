package limiting

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestSimpleParallelLimiting(t *testing.T) {
	stopCount := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopCount:
				return
			default:
				fmt.Printf("goroutines count: %v\n", runtime.NumGoroutine())
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg := new(sync.WaitGroup)
	limit := make(chan struct{}, 2)

	task := func(num int) {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		fmt.Printf("%v done\n", num)
		<-limit
	}

	const taskCount = 10
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		limit <- struct{}{}
		go task(i)
	}

	wg.Wait()
	close(stopCount)
}
