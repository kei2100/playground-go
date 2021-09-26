package limiting

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type goroutinesCounter struct {
	done chan struct{}
}

func newGoroutinesCounter() *goroutinesCounter {
	return &goroutinesCounter{
		done: make(chan struct{}),
	}
}

func (c *goroutinesCounter) start() <-chan int {
	count := make(chan int)
	go func() {
		defer close(count)
		for {
			select {
			case <-c.done:
				close(c.done)
				return
			case <-time.After(1 * time.Second):
				count <- runtime.NumGoroutine()
			}
		}
	}()
	return count
}

func (c *goroutinesCounter) stop() {
	c.done <- struct{}{}
	for range c.done {
		// wait for close(c.done)
	}
}

func TestSimpleParallelLimiting(t *testing.T) {
	counter := newGoroutinesCounter()
	countCh := counter.start()
	defer counter.stop()

	go func() {
		for c := range countCh {
			fmt.Printf("goroutines count: %v\n", c)
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
}
