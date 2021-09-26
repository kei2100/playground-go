package workerpool

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func worker(tasks <-chan func(), done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			t := <-tasks
			t()
		}
	}
}

func TestWorkerPool(t *testing.T) {
	const n = 10

	var wg sync.WaitGroup
	wg.Add(n)

	tasks := make(chan func(), n)
	done := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(tasks, done)
	}

	results := make([]int, n)
	for i := 0; i < n; i++ {
		num := i
		tasks <- func() {
			defer wg.Done()
			results[num] = num * num
		}
	}

	wg.Wait()

	fmt.Println(results)
}
