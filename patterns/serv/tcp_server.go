package serv

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var ErrServerClosed = errors.New("serv: server closed")

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

type TCPServerStats struct {
	NumConnections int
}

// base on http stdpkg
type TCPServer struct {
	mu          sync.Mutex
	state       atomic.Value
	ln          net.Listener
	connOpts    TCPConnOptions
	connTracker tcpConnTracker
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

		c := s.trackConn(conn)
		c.setOptions(s.connOpts)
		go handler(c)
	}
}

func (s *TCPServer) Close() error {
	err := s.CloseListener()
	for _, conn := range s.connTracker.all() {
		conn.Close()
	}
	return err
}

func (s *TCPServer) CloseListener() error {
	return s.withLockDo(func() error {
		if s.IsClosed() {
			return nil
		}
		err := s.ln.Close()
		s.ln = nil
		s.state.Store(stateClosed)
		if err != nil {
			return fmt.Errorf("serv: an error occurred when closing the listener: %v", err)
		}
		return nil
	})
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

const (
	stateClosed = iota
	stateListening
)

func (s *TCPServer) IsListening() bool {
	switch ss := s.state.Load().(type) {
	case int:
		return ss == stateListening
	case nil:
		// The state has not been set yet(= not listening)
		return false
	default:
		panic(fmt.Errorf("serv: invalid TCPServer.state type: %T", ss))
	}
}

func (s *TCPServer) IsClosed() bool {
	switch ss := s.state.Load().(type) {
	case int:
		return ss == stateClosed
	case nil:
		// The state has not been set yet(= closed)
		return true
	default:
		panic(fmt.Errorf("serv: invalid TCPServer.state type: %T", ss))
	}
}

func (s *TCPServer) setListener(ln net.Listener) error {
	return s.withLockDo(func() error {
		if s.IsListening() {
			return fmt.Errorf("serv: already serving")
		}
		s.ln = ln
		s.state.Store(stateListening)
		return nil
	})
}

func (s *TCPServer) Stats() TCPServerStats {
	return TCPServerStats{
		NumConnections: s.connTracker.count(),
	}
}

func (s *TCPServer) withLockDo(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

func (s *TCPServer) trackConn(conn net.Conn) *tcpServerConn {
	c := &tcpServerConn{s: s, Conn: conn}
	s.connTracker.add(c)
	return c
}

func (s *TCPServer) untrackConn(conn *tcpServerConn) {
	s.connTracker.remove(conn)
}

type tcpConnTracker struct {
	mu    sync.Mutex
	conns map[*tcpServerConn]struct{}
}

func (t *tcpConnTracker) add(conn *tcpServerConn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conns == nil {
		t.conns = make(map[*tcpServerConn]struct{})
	}
	t.conns[conn] = struct{}{}
}

func (t *tcpConnTracker) remove(conn *tcpServerConn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.conns, conn)
}

func (t *tcpConnTracker) all() []*tcpServerConn {
	t.mu.Lock()
	defer t.mu.Unlock()
	ret := make([]*tcpServerConn, len(t.conns))
	var i int
	for k := range t.conns {
		ret[i] = k
		i++
	}
	return ret
}

func (t *tcpConnTracker) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.conns)
}

type tcpServerConn struct {
	s *TCPServer
	net.Conn
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

func (c *tcpServerConn) setOptions(o TCPConnOptions) {
	if conn, ok := c.Conn.(*net.TCPConn); ok && o.KeepAlivePeriod > 0 {
		conn.SetKeepAlive(o.KeepAlive)
		conn.SetKeepAlivePeriod(o.KeepAlivePeriod)
	}
	if t := o.ReadTimeout; t > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(t))
	}
	if t := o.WriteTimeout; t > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(t))
	}
}

func (c *tcpServerConn) Close() error {
	c.s.untrackConn(c)
	return c.Conn.Close()
}
