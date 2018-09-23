package channel

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
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

func TestDynamicSelection(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	cases := make([]reflect.SelectCase, 3)
	for i, ch := range []interface{}{ch1, ch2, ch3} {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if chosen, recv, ok := reflect.Select(cases); ok {
				log.Printf("recv %v, ch%d chosen", recv, chosen+1)
				return
			}
		}
	}()

	r := rand.New(rand.NewSource(int64(os.Getpid())))
	i := r.Intn(3)
	ch := []chan int{ch1, ch2, ch3}[i]

	go func() { ch <- r.Intn(100) }()
	select {
	case <-done:
		break
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}
