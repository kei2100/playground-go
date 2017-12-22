package channel

import (
	"testing"
	"time"
)

func TestSendingTimeout(t *testing.T) {
	ch := make(chan struct{})

	f := func(msg chan<- string) {
		defer close(msg)
		select {
		// channelへの送信完了もselect caseすることができる
		case ch <- struct{}{}:
			msg <- "complete"
		case <-time.After(1 * time.Second):
			msg <- "timeout"
		}
	}

	msg := make(chan string)
	go f(msg)

	for m := range msg {
		if g, w := m, "timeout"; g != w {
			t.Errorf(" got %v, want %v", g, w)
		}
	}

	msg = make(chan string)
	go f(msg)
	<-ch

	for m := range msg {
		if g, w := m, "complete"; g != w {
			t.Errorf(" got %v, want %v", g, w)
		}
	}
}
