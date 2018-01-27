package startorder

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const n = 10

func TestEndOrder(t *testing.T) {
	// 終了順
	// channelに書き込んだ順に処理
	for m := range genEndOrder() {
		fmt.Println(m)
	}
}

func TestStartOrder(t *testing.T) {
	// 開始順
	// channelにchannelを書き込んだ順に処理
	ch := make(chan chan string, n)
	genStartOrder(ch)
	for ch := range ch {
		fmt.Println(<-ch)
	}
}

func genEndOrder() <-chan string {
	msg := make(chan string)

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(n-i) * time.Second / 2)
			msg <- fmt.Sprintf("num:%v done", i)
		}()
	}

	go func() {
		wg.Wait()
		close(msg)
	}()

	return msg
}

func genStartOrder(ch chan<- chan string) {
	defer close(ch)

	for i := 0; i < n; i++ {
		i := i
		msg := make(chan string)
		ch <- msg
		go func() {
			time.Sleep(time.Duration(n-i) * time.Second / 2)
			msg <- fmt.Sprintf("num:%v done", i)
		}()
	}
}
