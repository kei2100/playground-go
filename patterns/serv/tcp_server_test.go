package serv

import (
	"bufio"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestTCPServer_ServeClose(t *testing.T) {
	s := &TCPServer{
		handler: func(conn net.Conn) {
			defer conn.Close()
			// echo
			r := bufio.NewReader(conn)
			b, err := r.ReadBytes('\n')
			if err != nil {
				t.Fatal(err)
			}
			if _, err := conn.Write(b); err != nil {
				t.Fatal(err)
			}
		},
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	closed := make(chan struct{})
	go func() {
		err := s.Serve(ln)
		if g, w := err, ErrServerClosed; g != w {
			t.Errorf("Serve returns got %v, want %v", g, w)
		}
		close(closed)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello\n")); err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(b), "hello\n"; g != w {
		t.Errorf("received msg got %v, want %v", g, w)
	}

	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
		// test ok
	case <-time.After(3 * time.Second):
		t.Errorf("timeout exceeded while waiting for serv Close")
	}
}
