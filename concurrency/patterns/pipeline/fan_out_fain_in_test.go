package pipeline

import (
	"testing"
	"unicode/utf8"
	"sync"
	"fmt"
	"time"
	"context"
	"log"
)

func TestFanOutFanIn(t *testing.T) {
	// Fan-out: 複数のgoroutineを使い、同一のchannelから、それがcloseされるまで値を読み込むパターン

	// Fan-in: 複数のchannelを一つに束ねて、全てのchannelがcloseされるまで、束ねた一つのchannelから読み込みを行うパターン

	ctx, can := context.WithTimeout(context.Background(), 100 * time.Millisecond)
	defer can()

	in := gen(ctx, "fan-out", "fan-in", "test")

	// Fan-out
	o1 := wordCount(ctx, in)
	o2 := wordCount(ctx, in)
	o3 := wordCount(ctx, in)

	// Fan-in
	out := mergeChannels(ctx, o1, o2, o3)

	var count = 0
	for o := range out {
		count += o
		//// このsleepを付けるとctxがタイムアウトする
		//time.Sleep(110 * time.Millisecond)
	}

	fmt.Println(count)
}

func wordCount(ctx context.Context, in <-chan string) <-chan int {
	out := make(chan int)

	go func() {
		defer func() {
			log.Println("wordCount: finished. close out channel")
			close(out)
		}()

		for i := range in {
			select {
			case out <- utf8.RuneCountInString(i):
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Printf("wordCount: context was canceled or deadline exceeded %v", err.Error())
				}
				return
			}
		}
	}()

	return out
}

func mergeChannels(ctx context.Context, channels ...<-chan int) <-chan int{
	out := make(chan int)
	wg := new(sync.WaitGroup)

	reader := func(in <-chan int) {
		defer func() {
			log.Println("reader: finished. call wg.Done()")
			wg.Done()
		}()
		for i := range in {
			select {
			case out <- i:
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Printf("reader: context was canceled or deadline exceeded %v", err.Error())
				}
				return
			}
		}
	}

	for _, in := range channels {
		wg.Add(1)
		go reader(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
