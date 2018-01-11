package serv

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/kei2100/playground-go/util/once"
)

var ErrServerClosed = errors.New("serv: server closed")

// TODO set handler
// TODO avoid dup serve & close
// TODO configure timeout
// TODO shutdown

// base on http stdpkg
type TCPServer struct {
	lnCloser *once.Closer

	done    chan struct{}
	handler func(net.Conn)
}

func (s *TCPServer) Serve(ln net.Listener) error {
	s.lnCloser = once.NewCloser(ln)
	defer s.lnCloser.Close()

	s.done = make(chan struct{})
	var tempDelay time.Duration

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return ErrServerClosed
			default:
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
	close(s.done)
	return s.lnCloser.Close()
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}
