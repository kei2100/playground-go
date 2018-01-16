package serv

import (
	"context"
	"errors"
	"log"
	"net"
	"sync/atomic"
	"time"
)

var ErrServerClosed = errors.New("serv: server closed")

const notClosed int32 = 0
const closed int32 = 1

type listener struct {
	ln     net.Listener
	closed int32
}

func newListener(ln net.Listener) *listener {
	return &listener{ln: ln}
}

func (ln *listener) Close() error {
	if !atomic.CompareAndSwapInt32(&ln.closed, notClosed, closed) {
		return errors.New("serv: listener already closed")
	}
	return ln.ln.Close()
}

func (ln *listener) Closed() bool {
	c := atomic.LoadInt32(&ln.closed)
	return c == closed
}

// TODO set handler
// TODO avoid dup serve & close
// TODO configure timeout
// TODO shutdown

// base on http stdpkg
type TCPServer struct {
	ln      *listener
	handler func(net.Conn)
}

func (s *TCPServer) Serve(ln net.Listener) error {
	s.ln = newListener(ln)
	defer s.ln.Close()

	var tempDelay time.Duration

	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.ln.Closed() {
				return ErrServerClosed
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("serv: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0
		go s.handler(conn)
	}
}

func (s *TCPServer) Close() error {
	return s.ln.Close()
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}
