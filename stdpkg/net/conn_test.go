package net

import (
	"net"
	"sync"
	"testing"
)

func TestConnDoubleClose(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	serverConnCh, clientConnCh := make(chan net.Conn), make(chan net.Conn)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		serverConnCh <- conn
	}()

	go func() {
		conn, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		clientConnCh <- conn
	}()

	done := make(chan struct{})
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		for {
			select {
			case sc := <-serverConnCh:
				if err := sc.Close(); err != nil {
					t.Error(err)
				}
				// test double close
				// err e.g: close tcp 127.0.0.1:49762->127.0.0.1:49763: use of closed network connection
				if err := sc.Close(); err == nil {
					t.Error("got nil, want err")
				}
				wg.Done()
			case cc := <-clientConnCh:
				if err := cc.Close(); err != nil {
					t.Error(err)
				}
				// test double close
				// err e.g: close tcp 127.0.0.1:49763->127.0.0.1:49762: use of closed network connection
				if err := cc.Close(); err == nil {
					t.Error("got nil, want err")
				}
				wg.Done()
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	close(done)
}
