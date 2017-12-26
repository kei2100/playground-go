package proxy

import (
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	ln, err := Listen()
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error)
	go func() {
		err := ln.Serve(func(conn *net.TCPConn) {
			defer conn.Close()
			io.Copy(conn, conn)
		})
		errCh <- err
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// test echo
	send := []byte("echo")
	recv := make([]byte, len(send))

	conn.Write(send)
	conn.Read(recv)
	conn.Close()

	if !reflect.DeepEqual(send, recv) {
		t.Errorf("recv got %v, want %v", string(recv), string(send))
	}

	// test end of serve
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-errCh:
		if err == nil {
			t.Errorf("serve returns nil, want not nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout exceeded while waiting for the end of serve")
	}
}
