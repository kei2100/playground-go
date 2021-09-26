package net

import (
	"net"
	"testing"
)

func listenTCP(t *testing.T) *net.TCPListener {
	t.Helper()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}
	return ln
}

func serveTCP(t *testing.T, ln *net.TCPListener, handler func(*net.TCPConn)) error {
	t.Helper()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		go handler(conn)
	}
}

func dialTCP(t *testing.T, addr net.Addr, callback func(*net.TCPConn)) {
	t.Helper()

	taddr, err := net.ResolveTCPAddr("tcp", addr.String())
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}

	conn, err := net.DialTCP("tcp", nil, taddr)
	if err != nil {
		t.Errorf("failed to dial: %v", err)
	}
	callback(conn)
}
