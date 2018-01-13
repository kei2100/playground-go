package net

import (
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestConnDoubleClose(t *testing.T) {
	t.Parallel()

	ln := listenTCP(t)
	defer ln.Close()

	srvdone := make(chan struct{})

	go serveTCP(t, ln, func(conn *net.TCPConn) {
		defer close(srvdone)

		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		// test double close
		// err e.g: close tcp 127.0.0.1:49762->127.0.0.1:49763: use of closed network connection
		if err := conn.Close(); err == nil {
			t.Error("got nil, want err")
		}
	})

	dialTCP(t, ln.Addr(), func(conn *net.TCPConn) {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		// test double close
		// err e.g: close tcp 127.0.0.1:49763->127.0.0.1:49762: use of closed network connection
		if err := conn.Close(); err == nil {
			t.Error("got nil, want err")
		}
	})

	select {
	case <-srvdone:
		return
	case <-time.After(3 * time.Second):
		t.Error("timeout exceeded while waiting for servdone")
	}
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
