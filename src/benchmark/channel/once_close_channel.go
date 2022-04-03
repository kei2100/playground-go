package channel

import "sync"

// ChanHolder1 using select case
type ChanHolder1 struct {
	ch chan struct{}
}

func NewChanHolder1() *ChanHolder1 {
	return &ChanHolder1{ch: make(chan struct{})}
}

func (c *ChanHolder1) OnceCloseChannel() {
	select {
	case <-c.ch:
		return
	default:
		close(c.ch)
	}
}

// ChanHolder2 using sync.Once
type ChanHolder2 struct {
	ch   chan struct{}
	once sync.Once
}

func NewChanHolder2() *ChanHolder2 {
	return &ChanHolder2{ch: make(chan struct{})}
}

func (c *ChanHolder2) OnceCloseChannel() {
	c.once.Do(func() {
		close(c.ch)
	})
}
