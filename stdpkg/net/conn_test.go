package net

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"runtime"
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

func TestConnSetKeepalive(t *testing.T) {
	kaln := listenTCP(t)
	nokaln := listenTCP(t)

	fmt.Printf("keepalive server addr   : %v\n", kaln.Addr())
	fmt.Printf("no keepalive server addr: %v\n", nokaln.Addr())

	go func() {
		serveTCP(t, kaln, func(conn *net.TCPConn) {
			defer conn.Close()
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(time.Second)

			r := bufio.NewReader(conn)
			b, err := r.ReadBytes('\n')
			fmt.Printf("keep alive server read finished: n bytes %v: err %v\n", len(b), err)
		})
	}()

	go func() {
		serveTCP(t, nokaln, func(conn *net.TCPConn) {
			defer conn.Close()

			r := bufio.NewReader(conn)
			b, err := r.ReadBytes('\n')
			fmt.Printf("no keep alive server read finished: n bytes %v: err %v\n", len(b), err)
		})
	}()

	go dialTCP(t, kaln.Addr(), func(conn *net.TCPConn) {
		defer conn.Close()
		fmt.Printf("receive keepalive client addr   :%v\n", conn.LocalAddr())

		r := bufio.NewReader(conn)
		b, err := r.ReadBytes('\n')
		fmt.Printf("receive keepalive client read finished: n bytes %v: err %v\n", len(b), err)
	})
	go dialTCP(t, nokaln.Addr(), func(conn *net.TCPConn) {
		defer conn.Close()
		fmt.Printf("receive no keepalive client addr:%v\n", conn.LocalAddr())

		r := bufio.NewReader(conn)
		b, err := r.ReadBytes('\n')
		fmt.Printf("receive no keepalive client read finished: n bytes %v: err %v\n", len(b), err)
	})

	// wait to TERMINATE
	for {
		fmt.Printf("num goroutines: %v\n", runtime.NumGoroutine())
		time.Sleep(3 * time.Second)
	}

	// (a) tcpdump -i lo0 src port {keepalive server port}
	// (b) tcpdump -i lo0 src port {no keepalive server port}
	// すると、
	// (a)側には1秒ごとにkeepalive probeパケットがクライアントに送信されるのが確認できる。
	// (b)は何も流れない。

	// OSXで /etc/pf.confに、
	// block in on lo0 proto tcp from self port {ka server port} to self port {ka client port}
	// を追記して、
	// pfctl -f /etc/pf.conf （pf.confの内容反映）
	// pfctl -e（pfを有効化）
	// すると、
	// keepalive probeがclientに届かなくなるため、
	// probeカウント超過後（デフォ9回？）にサーバー側でconn.readがerrで終了することを確認できる。
	// カウント超過前に pfctl -d（pfを無効化）すれば、errにはならずprobe送信を継続する。
}
