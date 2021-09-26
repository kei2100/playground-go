// +build darwin freebsd linux

package main

import (
	"log"
	"os"
	"time"

	"bitbucket.org/avd/go-ipc/fifo"
)

func main() {
	go func() {
		wfifo, err := fifo.New("fifo", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("new: %v", err)
		}
		wfifo.Close()
		log.Println(wfifo.Write([]byte{1}))
		defer wfifo.Close()
	}()
	//buff := make([]byte, 8)
	rfifo, err := fifo.New("fifo", os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		panic("new")
	}
	rfifo.Close()
	time.Sleep(time.Second)
	//defer rfifo.Close()
	//n, err := rfifo.Read(buff)
	//log.Printf("i: %d, err : %v, io.EOF: %v. %T", n, err, err == io.EOF, err)
}
