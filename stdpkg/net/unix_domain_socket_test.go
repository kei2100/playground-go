package net

import (
	"path/filepath"
	"os"
	"net"
	"fmt"
	"io"
	"testing"
)

// TODO server pattern
type unixSocketServer struct {
	sockpath string
	stopSignal chan struct{}
}

func (s *unixSocketServer) listenAndServe() error {
	s.sockpath = filepath.Join(os.TempDir(), "unix_socket_server_test.sock")
	ln, err := net.Listen("unix", s.sockpath)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer func() {
		if err := ln.Close(); err != nil {
			fmt.Printf("listener close error :%v", err)
		}
		if err := os.Remove(s.sockpath); err != nil {
			fmt.Printf("faield to remove %v :%v", s.sockpath, err)
		}
	}()

	handler := func(conn net.Conn) {
		defer conn.Close()
		io.Copy(conn, conn)
	}
	acceptor := func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("failed to accept: %v", err)
				return
			}
			go handler(conn)
		}
	}
	go acceptor()

	s.stopSignal = make(chan struct{})
	<-s.stopSignal

	return fmt.Errorf("server stopSignal received")
}

func (s *unixSocketServer) stop() {
	close(s.stopSignal)
}

func TestUnixSocketServer(t *testing.T) {

}
