package channel

import (
	"sync"
	"fmt"
)

func ExampleUnidirectionalChannel() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	sender := func(c chan<- int) {
		c <- 1
		wg.Done()
	}
	receiver := func(c <-chan int) {
		v := <-c
		fmt.Println(v)
		wg.Done()

		// Output:
		// 1
	}

	wg.Add(1)
	go sender(ch)
	wg.Add(1)
	go receiver(ch)

	wg.Wait()
}
