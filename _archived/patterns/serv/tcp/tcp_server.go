package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ErrServerClosed is returned by the Server's Serve method
// after a call to Close.
var ErrServerClosed = errors.New("serv: server closed")

// ServerOptions is option for Server
type ServerOptions func(*Server)

// WithKeepAlive set keepalive to incoming connection
func WithKeepAlive(keepalive bool) ServerOptions {
	return func(s *Server) {
		s.connOpts.KeepAlive = keepalive
	}
}

// WithKeepAlivePeriod set keepalive period to incoming connection
func WithKeepAlivePeriod(p time.Duration) ServerOptions {
	return func(s *Server) {
		s.connOpts.KeepAlivePeriod = p
	}
}

// WithReadTimeout set read timeout to incoming connection
func WithReadTimeout(t time.Duration) ServerOptions {
	return func(s *Server) {
		s.connOpts.ReadTimeout = t
	}
}

// WithWriteTimeout set write timeout to incoming connection
func WithWriteTimeout(t time.Duration) ServerOptions {
	return func(s *Server) {
		s.connOpts.WriteTimeout = t
	}
}

// WithConnOptions set ConnOptions to incoming connection
func WithConnOptions(o ConnOptions) ServerOptions {
	return func(s *Server) {
		s.connOpts = o
	}
}

// HandleFunc is type of handler
type HandleFunc func(net.Conn)

// ServerStats is statistics of the Server
type ServerStats struct {
	NumConnections int
}

// A Server defines parameters for running an TCP server.
// The zero value for Server is a valid configuration.
type Server struct {
	mu          sync.Mutex
	state       atomic.Value
	ln          net.Listener
	connOpts    ConnOptions
	connTracker connTracker
}

// Serve accepts incoming TCP connections on the listener ln,
// creating a new service goroutine for each. The service goroutines call handler to reply to them.
func (s *Server) Serve(ln net.Listener, handler HandleFunc, opts ...ServerOptions) error {
	s.setOptions(opts...)
	if err := s.setListener(ln); err != nil {
		return err
	}

	var tempDelay time.Duration

	for {
		conn, err := ln.Accept()
		if err != nil {
			if s.IsClosing() || s.IsClosed() {
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

// Close immediately closes the server listener and all pending connections
func (s *Server) Close() error {
	err := s.CloseListener()
	for _, conn := range s.connTracker.all() {
		conn.Close()
	}
	return err
}

// CloseListener closes the server listener. It stops accepting new connections
func (s *Server) CloseListener() error {
	return s.withLockDo(func() error {
		if s.IsClosing() || s.IsClosed() {
			return nil
		}
		s.state.Store(stateClosing)
		err := s.ln.Close()
		s.ln = nil
		s.state.Store(stateClosed)
		if err != nil {
			return fmt.Errorf("serv: an error occurred when closing the listener: %v", err)
		}
		return nil
	})
}

func (s *Server) setOptions(opts ...ServerOptions) {
	s.withLockDo(func() error {
		s.connOpts = defaultConnOptions()
		for _, o := range opts {
			o(s)
		}
		return nil
	})
}

const (
	stateClosed = iota
	stateClosing
	stateListening
)

// IsClosed reports whether the server listener is closed
func (s *Server) IsClosed() bool {
	switch ss := s.state.Load().(type) {
	case int:
		return ss == stateClosed
	case nil:
		// The state has not been set yet(= closed)
		return true
	default:
		panic(fmt.Errorf("serv: invalid Server.state type: %T", ss))
	}
}

// IsClosing reports whether the server listener is Closing
func (s *Server) IsClosing() bool {
	switch ss := s.state.Load().(type) {
	case int:
		return ss == stateClosing
	case nil:
		// The state has not been set yet(= not closing)
		return false
	default:
		panic(fmt.Errorf("serv: invalid Server.state type: %T", ss))
	}
}

// IsListening reports whether the server listener is listening
func (s *Server) IsListening() bool {
	switch ss := s.state.Load().(type) {
	case int:
		return ss == stateListening
	case nil:
		// The state has not been set yet(= not listening)
		return false
	default:
		panic(fmt.Errorf("serv: invalid Server.state type: %T", ss))
	}
}

func (s *Server) setListener(ln net.Listener) error {
	return s.withLockDo(func() error {
		if s.IsListening() {
			return fmt.Errorf("serv: already serving")
		}
		s.ln = ln
		s.state.Store(stateListening)
		return nil
	})
}

// Stats return ServerStats
func (s *Server) Stats() ServerStats {
	return ServerStats{
		NumConnections: s.connTracker.count(),
	}
}

func (s *Server) withLockDo(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

func (s *Server) trackConn(conn net.Conn) *serverConn {
	c := &serverConn{s: s, Conn: conn}
	s.connTracker.add(c)
	return c
}

func (s *Server) untrackConn(conn *serverConn) {
	s.connTracker.remove(conn)
}

type connTracker struct {
	mu    sync.Mutex
	conns map[*serverConn]struct{}
}

func (t *connTracker) add(conn *serverConn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conns == nil {
		t.conns = make(map[*serverConn]struct{})
	}
	t.conns[conn] = struct{}{}
}

func (t *connTracker) remove(conn *serverConn) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.conns, conn)
}

func (t *connTracker) all() []*serverConn {
	t.mu.Lock()
	defer t.mu.Unlock()
	ret := make([]*serverConn, len(t.conns))
	var i int
	for k := range t.conns {
		ret[i] = k
		i++
	}
	return ret
}

func (t *connTracker) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.conns)
}

type serverConn struct {
	s *Server
	net.Conn
}

// ConnOptions defines parameter for TCP connection
type ConnOptions struct {
	KeepAlive       bool
	KeepAlivePeriod time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func defaultConnOptions() ConnOptions {
	return ConnOptions{
		KeepAlive:       true,
		KeepAlivePeriod: 1 * time.Minute,
		ReadTimeout:     3 * time.Minute,
	}
}

func (c *serverConn) setOptions(o ConnOptions) {
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

// Close closes the connection and turn off server tracking.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *serverConn) Close() error {
	c.s.untrackConn(c)
	return c.Conn.Close()
}
