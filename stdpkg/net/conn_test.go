package net

import (
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func TestConnDoubleClose(t *testing.T) {
	t.Parallel()

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

func listenTCP(t *testing.T) *net.TCPListener {
	t.Helper()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}
	return ln
}

func serveTCP(t *testing.T, ln *net.TCPListener, handler func(*net.TCPConn)) error {
	t.Helper()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		go handler(conn)
	}
}

func dialTCP(t *testing.T, addr net.Addr, callback func(*net.TCPConn)) {
	t.Helper()

	taddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}

	conn, err := net.DialTCP("tcp", nil, taddr)
	if err != nil {
		t.Errorf("failed to dial: %v", err)
	}
	callback(conn)
}

func TestConnSetDeadline(t *testing.T) {
	t.Parallel()

	type result struct {
		err     error
		elapsed int64
	}
	rch := make(chan result)

	ln := listenTCP(t)
	go func() {
		serveTCP(t, ln, func(conn *net.TCPConn) {
			defer conn.Close()
			n := time.Now()
			conn.SetDeadline(time.Now().Add(1 * time.Second))

			_, err := ioutil.ReadAll(conn)
			rch <- result{
				err:     err,
				elapsed: time.Now().Unix() - n.Unix(),
			}
		})
	}()

	dialTCP(t, ln.Addr(), func(conn *net.TCPConn) {
		defer conn.Close()
		time.Sleep(2 * time.Second)
	})

	select {
	case r := <-rch:
		if r.err == nil {
			t.Error("result.error got nil, want an error")
		}
		if r.elapsed < 1 {
			t.Errorf("result.elapsed got %v, want greater than 1", r.elapsed)
		}
		err, ok := r.err.(net.Error)
		if !ok {
			t.Errorf("result.error got %T, want an net.Error", r.err)
			return
		}
		if !err.Timeout() || !err.Temporary() {
			t.Errorf("timeout %v, temporary %v, want both true", err.Timeout(), err.Temporary())
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout exceeded while waiting for rch")
	}
}
