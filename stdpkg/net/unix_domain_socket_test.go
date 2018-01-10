package net

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

type unixSocket struct {
	path     string
	listener net.Listener
}

func listenUnixSocket() (*unixSocket, error) {
	p := filepath.Join(os.TempDir(), "unix_socket_server_test.sock")
	os.Remove(p)

	ln, err := net.Listen("unix", p)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}
	return &unixSocket{path: p, listener: ln}, nil
}

type unixSocketServer struct {
	sock *unixSocket
	done chan struct{}
}

func newUnixSocketServer(sock *unixSocket) *unixSocketServer {
	return &unixSocketServer{sock: sock, done: make(chan struct{})}
}

func (s *unixSocketServer) serve() error {
	defer func() {
		s.sock.listener.Close()
		os.Remove(s.sock.path)
	}()

	for {
		conn, err := s.sock.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return errors.New("server closed")
			default:
				return err
			}
		}
		go func() {
			defer conn.Close()
			r := bufio.NewReader(conn)
			b, err := r.ReadBytes('\n')
			if err != nil {
				fmt.Printf("failed to read %v", err)
			}
			conn.Write(b)
		}()
	}
}

func (s *unixSocketServer) close() {
	close(s.done)
	s.sock.listener.Close()
}

func TestUnixSocketServer(t *testing.T) {
	sock, err := listenUnixSocket()
	if err != nil {
		t.Fatal(err)
	}

	srv := newUnixSocketServer(sock)
	go srv.serve()
	defer srv.close()

	conn, err := net.Dial("unix", sock.path)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte("hello\n"))
	got, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(got), "hello\n"; g != w {
		t.Errorf("received got %v, want %v", g, w)
	}
}
