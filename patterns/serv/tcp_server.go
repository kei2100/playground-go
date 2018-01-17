package serv

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
	"fmt"
)

var ErrServerClosed = errors.New("serv: server closed")

const (
	stateClosed = iota
	stateListening
)

type servState int32

func (st *servState) IsListening() bool {
	return *st == stateListening
}

func (st *servState) setListening() {
	*st = stateListening
}

func (st *servState) IsClosed() bool {
	return *st == stateClosed
}

func (st *servState) setClosed() {
	*st = stateClosed
}

// TODO set handler
// TODO avoid dup serve & close
// TODO configure timeout
// TODO shutdown

// base on http stdpkg
type TCPServer struct {
	mu sync.Mutex
	servState
	ln      net.Listener
	handler func(net.Conn)
}

func (s *TCPServer) Serve(ln net.Listener) error {
	s.mu.Lock()
	if s.IsListening() {
		return fmt.Errorf("serv: already listening")
	}
	s.ln = ln
	s.setListening()
	s.mu.Unlock()

	var tempDelay time.Duration

	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.IsClosed() {
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
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsListening() {
		return nil
	}
	err := s.ln.Close()
	s.ln = nil
	s.setClosed()

	return err
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}
