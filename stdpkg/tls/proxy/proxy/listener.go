package proxy

import (
	"fmt"
	"net"
	"time"
)

// HandleFunc is the handler for accepted Connection
type HandleFunc func(conn *net.TCPConn)

// ConnOption is the functional option for accepted connection
type ConnOption func(conn *net.TCPConn)

// SetKeepAlive to TCP Connection
func SetKeepAlive(keepalive bool) ConnOption {
	return func(conn *net.TCPConn) {
		conn.SetKeepAlive(keepalive)
	}
}

// SetKeepAlivePeriod to TCP Connection
func SetKeepAlivePeriod(d time.Duration) ConnOption {
	return func(conn *net.TCPConn) {
		conn.SetKeepAlivePeriod(d)
	}
}

// SetNoDelay to TCP Connection
func SetNoDelay(noDelay bool) ConnOption {
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

// Serve calls handler to handle the incoming connections and continue waiting for the next connections.
// Serve always returns a non-nil error.
func (ln *Listener) Serve(handler HandleFunc, opts ...ConnOption) error {
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		for _, o := range opts {
			o(conn)
		}
		handler(conn)
	}
}
