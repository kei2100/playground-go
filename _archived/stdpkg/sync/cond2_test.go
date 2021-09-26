package sync

import (
	"log"
	"sync"
	"testing"
	"time"
)

// from ISBN978-4-87311-856-8
func TestCondSignal(t *testing.T) {
	c := sync.NewCond(&sync.Mutex{})
	const maxQLen = 10
	q := make([]int, 0, 10)

	remove := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		p := q[0]
		q = q[1:]
		log.Printf("removed: %d", p)
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < maxQLen; i++ {
		c.L.Lock()
		for len(q) == 2 {
			c.Wait()
		}
		q = append(q, i)
		log.Printf("append: %d", i)
		go remove(time.Millisecond)
		c.L.Unlock()
	}

	// === RUN   TestCondSignal
	// 2018/12/01 21:47:13 append: 0
	// 2018/12/01 21:47:13 append: 1
	// 2018/12/01 21:47:13 removed: 0
	// 2018/12/01 21:47:13 append: 2
	// 2018/12/01 21:47:13 removed: 1
	// 2018/12/01 21:47:13 append: 3
	// 2018/12/01 21:47:13 removed: 2
	// 2018/12/01 21:47:13 append: 4
	// 2018/12/01 21:47:13 removed: 3
	// 2018/12/01 21:47:13 append: 5
	// 2018/12/01 21:47:13 removed: 4
	// 2018/12/01 21:47:13 append: 6
	// 2018/12/01 21:47:13 removed: 5
	// 2018/12/01 21:47:13 append: 7
	// 2018/12/01 21:47:13 removed: 6
	// 2018/12/01 21:47:13 append: 8
	// 2018/12/01 21:47:13 removed: 7
	// 2018/12/01 21:47:13 append: 9
}
