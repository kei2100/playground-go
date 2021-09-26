package main

import (
	"os"
	"sync"
	"testing"

	"bitbucket.org/avd/go-ipc/mq"
)

func TestCloseAndRead(t *testing.T) {
	mq.Destroy("mq")
	q, err := mq.New("mq", os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		//q.Send([]byte{1})
	}()
	q.Close()
	go func() {
		defer wg.Done()
		//b := make([]byte, 1)
		//q.Receive(b)
	}()

	wg.Wait()
}
