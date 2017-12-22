package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	StateNotReady int32 = iota
	StateSticking
	StateHeld
)

type beachFlag struct {
	state int32
}

func (a *beachFlag) stick() {
	atomic.StoreInt32(&a.state, StateSticking)
}

func (a *beachFlag) notReady() bool {
	return atomic.LoadInt32(&a.state) == StateNotReady
}

func (a *beachFlag) hold() bool {
	return atomic.CompareAndSwapInt32(&a.state, StateSticking, StateHeld)
}

func TestCond(t *testing.T) {
	m := new(sync.Mutex)
	c := sync.NewCond(m)
	done := new(sync.WaitGroup)

	f := new(beachFlag)

	gr := func(num int) {
		defer done.Done()
		c.L.Lock()
		if f.notReady() {
			fmt.Printf("goroutine(%v) wait\n", num)
			c.Wait()
		}
		c.L.Unlock()

		if f.hold() {
			fmt.Printf("goroutine(%v) hold successful\n", num)
		} else {
			fmt.Printf("goroutine(%v) hold failure\n", num)
		}
	}

	const grSize = 20
	done.Add(grSize)
	for i := 0; i < grSize; i++ {
		go gr(i)
	}

	c.L.Lock()
	f.stick()
	c.L.Unlock()
	c.Broadcast()

	done.Wait()
}
