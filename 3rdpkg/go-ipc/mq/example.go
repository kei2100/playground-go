package main

import (
	"log"
	"os"
	"time"

	"bitbucket.org/avd/go-ipc/mq"
)

func main() {
	mq.Destroy("mq")
	q, err := mq.New("mq", os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatalf("new queue: %v", err)
	}
	defer q.Close()
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	go func() {
		if err := q.Send(data); err != nil {
			panic("send")
		}
		log.Println("done send")
		recv2 := make([]byte, len(data))
		q.Receive(recv2)
		log.Println(recv2)
		log.Println("done recv2")
	}()

	q2, err := mq.Open("mq", 0)
	if err != nil {
		panic("open")
	}
	defer q2.Close()
	received := make([]byte, len(data))
	l, err := q2.Receive(received)
	if err != nil {
		panic("receive")
	}
	if l != len(data) {
		panic("wrong len")
	}
	for i, b := range received {
		if b != data[i] {
			panic("wrong data")
		}
	}
	log.Println("done receive")

	log.Println("send 2")
	//if err := q2.Send([]byte{1,1}); err != nil {
	if err := q2.Send([]byte{8, 8, 8, 8, 8, 8, 8, 8}); err != nil {
		log.Fatalf("send2: %v", err)
	}
	time.Sleep(time.Second)
}
