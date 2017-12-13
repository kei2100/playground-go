package pipeline

import (
	"testing"
	"unicode/utf8"
	"sync"
	"fmt"
)

func TestFanOutFanIn(t *testing.T) {
	// Fan-out: 複数のgoroutineを使い、同一のchannelから、それがcloseされるまで値を読み込むパターン

	// Fan-in: 複数のchannelを一つに束ねて、全てのchannelがcloseされるまで、束ねた一つのchannelから読み込みを行うパターン

	in := gen("fan-out", "fan-in", "test")

	// Fan-out
	o1 := wordCount(in)
	o2 := wordCount(in)
	o3 := wordCount(in)

	// Fan-in
	out := mergeChannels(o1, o2, o3)

	var count = 0
	for o := range out {
		count += o
	}

	fmt.Println(count)
}

func wordCount(in <-chan string) <-chan int {
	out := make(chan int)

	go func() {
		for i := range in {
			out <- utf8.RuneCountInString(i)
		}
		close(out)
	}()

	return out
}

func mergeChannels(channels ...<-chan int) <-chan int{
	out := make(chan int)
	wg := new(sync.WaitGroup)

	reader := func(in <-chan int) {
		for i := range in {
			out <- i
		}
		wg.Done()
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
