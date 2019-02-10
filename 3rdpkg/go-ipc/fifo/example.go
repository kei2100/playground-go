package main

import (
	"log"
	"os"

	"bitbucket.org/avd/go-ipc/fifo"
)

func main() {
	testData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	go func() {
		fifo, err := fifo.New("fifo", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("new: %v", err)
		}
		defer fifo.Close()
		defer fifo.Destroy()
		if written, err := fifo.Write(testData); err != nil || written != len(testData) {
			panic("write")
		}
	}()
	buff := make([]byte, len(testData))
	fifo, err := fifo.New("fifo", os.O_CREATE|os.O_RDONLY, 0666)
	//fifo, err := fifo.New("fifo", os.O_CREATE|os.O_RDONLY|fifo.O_NONBLOCK, 0666) // non blocking
	if err != nil {
		panic("new")
	}
	defer fifo.Close()
	if read, err := fifo.Read(buff); err != nil || read != len(testData) {
		panic("read")
	}
	// ensure we've received valid data
	for i, b := range buff {
		println(b)
		if b != testData[i] {
			panic("wrong data")
		}
	}
}
