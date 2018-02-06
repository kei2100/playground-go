package gateway

import (
	"fmt"
	"net"
	"time"
)

// RouteFunc is the router for accepted Connection
type RouteFunc func(conn net.Conn)

// ConnOption is the functional option for accepted connection
type ConnOption func(conn *net.TCPConn)

// WithKeepAlive option
func WithKeepAlive(keepalive bool) ConnOption {
	return func(conn *net.TCPConn) {
		conn.SetKeepAlive(keepalive)
	}
}

// WithKeepAlivePeriod option
func WithKeepAlivePeriod(d time.Duration) ConnOption {
	return func(conn *net.TCPConn) {
		conn.SetKeepAlivePeriod(d)
	}
}

// WithNoDelay option
func WithNoDelay(noDelay bool) ConnOption {
	return func(conn *net.TCPConn) {
		conn.SetNoDelay(noDelay)
	}
}

// Listener for the proxy
type Listener struct {
	*net.TCPListener
}

// Listen on a free TCP port on localhost
func Listen() (*Listener, error) {
	const tcp = "tcp"

	addr, err := net.ResolveTCPAddr(tcp, "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("proxy: failed to resolve addr: %v", err)
	}

	ln, err := net.ListenTCP(tcp, addr)
	if err != nil {
		return nil, fmt.Errorf("proxy: failed to listen %v: %v", ln.Addr(), err)
	}
	return &Listener{TCPListener: ln}, nil
}

// Serve calls router to handle the incoming connections and continue waiting for the next connections.
// Serve always returns a non-nil error.
func (ln *Listener) Serve(router RouteFunc, opts ...ConnOption) error {
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		for _, o := range opts {
			o(conn)
		}
		go router(conn)
	}
}
