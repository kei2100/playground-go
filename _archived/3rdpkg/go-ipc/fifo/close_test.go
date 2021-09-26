package main_test

import (
	"log"
	"os"
	"sync"
	"testing"

	"bitbucket.org/avd/go-ipc/fifo"
)

func TestFIFO(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		wfifo, err := fifo.New("fifo-test", os.O_CREATE|os.O_WRONLY, 0666)
		panicIf(err)
		defer wfifo.Close()

		wfifo.Write(data)
	}()

	rfifo, err := fifo.New("fifo-test", os.O_CREATE|os.O_RDONLY, 0666)
	panicIf(err)
	defer rfifo.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		b := make([]byte, len(data))
		_, err := rfifo.Read(b)
		log.Println(err)
	}()

	wg.Wait()
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
