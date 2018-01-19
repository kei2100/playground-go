package serv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
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

type TCPConnOptions struct {
	KeepAlive       bool
	KeepAlivePeriod time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func defaultTCPConnOptions() TCPConnOptions {
	return TCPConnOptions{
		KeepAlive:       true,
		KeepAlivePeriod: 1 * time.Minute,
		ReadTimeout:     3 * time.Minute,
	}
}

type TCPServerOptions func(*TCPServer)

func WithKeepAlive(keepalive bool) TCPServerOptions {
	return func(s *TCPServer) {
		s.connOpts.KeepAlive = keepalive
	}
}

func WithKeepAlivePeriod(p time.Duration) TCPServerOptions {
	return func(s *TCPServer) {
		s.connOpts.KeepAlivePeriod = p
	}
}

func WithReadTimeout(t time.Duration) TCPServerOptions {
	return func(s *TCPServer) {
		s.connOpts.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) TCPServerOptions {
	return func(s *TCPServer) {
		s.connOpts.WriteTimeout = t
	}
}

func WithConnOptions(o TCPConnOptions) TCPServerOptions {
	return func(s *TCPServer) {
		s.connOpts = o
	}
}

type TCPHandleFunc func(net.Conn)

// TODO configure timeout
// TODO shutdown
// TODO error msg

// base on http stdpkg
type TCPServer struct {
	mu sync.Mutex

	servState
	ln net.Listener

	connOpts TCPConnOptions
}

func (s *TCPServer) Serve(ln net.Listener, handler TCPHandleFunc, opts ...TCPServerOptions) error {
	s.setOptions(opts...)
	if err := s.setListener(ln); err != nil {
		return err
	}

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

		go handler(conn)
	}
}

func (s *TCPServer) Close() error {
	return s.closeListener()
}

func (s *TCPServer) Shutdown(ctx context.Context) error {
	return nil
}

func (s *TCPServer) setOptions(opts ...TCPServerOptions) {
	s.withLockDo(func() error {
		s.connOpts = defaultTCPConnOptions()
		for _, o := range opts {
			o(s)
		}
		return nil
	})
}

func (s *TCPServer) setListener(ln net.Listener) error {
	return s.withLockDo(func() error {
		if s.IsListening() {
			return fmt.Errorf("serv: already listening")
		}
		s.ln = ln
		s.setListening()
		return nil
	})
}

func (s *TCPServer) closeListener() error {
	return s.withLockDo(func() error {
		if s.IsClosed() {
			return nil
		}
		err := s.ln.Close()
		s.ln = nil
		s.setClosed()

		return err
	})
}

func (s *TCPServer) withLockDo(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

type tcpServerConn struct {
	net.Conn
	s *TCPServer
}
