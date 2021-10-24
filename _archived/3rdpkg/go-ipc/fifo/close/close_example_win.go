//go:build windows
// +build windows

package main

import (
	"log"
	"os"
	"syscall"

	"bitbucket.org/avd/go-ipc/fifo"
)

func main() {
	go func() {
		wfifo, err := fifo.New("fifo", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("new: %v", err)
		}
		defer wfifo.Close()
	}()
	buff := make([]byte, 8)
	rfifo, err := fifo.New("fifo", os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		panic("new")
	}
	defer rfifo.Close()
	n, err := rfifo.Read(buff)
	log.Printf("i: %d, err : %v, syscall.ERROR_BROKEN_PIPE: %v. %T", n, err, err == syscall.ERROR_BROKEN_PIPE, err)
}
