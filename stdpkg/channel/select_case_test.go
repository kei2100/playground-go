package channel

import (
	"fmt"
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

func TestForSelect(t *testing.T) {
	ch := make(chan int)

	go func() {
		i := 0
		for {
			select {
			case ch <- i:
				return
			default:
				// chが受信されなければdefaultにいき、次のループに入る
				time.Sleep(1 * time.Millisecond)
				i++
			}
		}
	}()

	time.Sleep(10 * time.Millisecond)
	fmt.Println(<-ch)
}
