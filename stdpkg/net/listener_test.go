package net

import (
	"fmt"
	"net"
	"testing"
	"time"
)

const tcp = "tcp"

func TestListenFreePort(t *testing.T) {
	// Linux では、ポート番号として 0 を指定して bind() を呼び出した場合や、
	// bind() を呼び出さずに connect() を呼び出した場合などに、
	// 未使用のローカルポート番号が自動的に割り当てられる。
	addr, err := net.ResolveTCPAddr(tcp, "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.ListenTCP(tcp, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	fmt.Println(ln.Addr())
}

func TestListenerAcceptFromClosedListener(t *testing.T) {
	addr, err := net.ResolveTCPAddr(tcp, "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.ListenTCP(tcp, addr)
	if err != nil {
		t.Fatal(err)
	}

	err = ln.Close()
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		if _, err = ln.Accept(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-errCh:
		// err = e.g. accept tcp 127.0.0.1:62302: use of closed network connection
		return
	case <-time.After(1 * time.Second):
		t.Errorf("timeout exceeded while waiting for send error")
	}
}

func TestListenerAcceptDeadlineExceeded(t *testing.T) {
	addr, err := net.ResolveTCPAddr(tcp, "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.ListenTCP(tcp, addr)
	if err != nil {
		t.Fatal(err)
	}

	err = ln.SetDeadline(time.Now().Add(10 * time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		if _, err = ln.Accept(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-errCh:
		// err = e.g. accept  tcp 127.0.0.1:62330: i/o timeout
		return
	case <-time.After(1 * time.Second):
		t.Errorf("timeout exceeded while waiting for send error")
	}
}
