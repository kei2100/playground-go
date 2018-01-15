package serv

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestTCPServer_ServeClose(t *testing.T) {
	t.Parallel()

	s := &TCPServer{
		handler: func(conn net.Conn) {
			defer conn.Close()
			// echo
			r := bufio.NewReader(conn)
			b, err := r.ReadBytes('\n')
			if err != nil {
				t.Fatal(err)
			}
			if _, err := conn.Write(b); err != nil {
				t.Fatal(err)
			}
		},
	}

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	closed := make(chan struct{})
	go func() {
		err := s.Serve(ln)
		if g, w := err, ErrServerClosed; g != w {
			t.Errorf("Serve returns got %v, want %v", g, w)
		}
		close(closed)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("hello\n")); err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(b), "hello\n"; g != w {
		t.Errorf("received msg got %v, want %v", g, w)
	}

	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
		// test ok
	case <-time.After(3 * time.Second):
		t.Errorf("timeout exceeded while waiting for serv Close")
	}
}

// implements net.Listener. Accept() always return err
type acceptErrorListener struct {
	raise error
}

func (l *acceptErrorListener) Accept() (net.Conn, error) {
	return nil, l.raise
}

func (l *acceptErrorListener) Close() error {
	return nil
}

func (l *acceptErrorListener) Addr() net.Addr {
	a, _ := net.ResolveTCPAddr("tcp", "192.0.2.0:0")
	return a
}

// represents a temporary net.Error
type temporaryError struct{}

func (e *temporaryError) Error() string {
	return "temporary error"
}

func (e *temporaryError) Timeout() bool {
	return true
}

func (e *temporaryError) Temporary() bool {
	return true
}

func TestTCPServer_HandleAcceptError(t *testing.T) {
	t.Parallel()

	t.Run("temporary error", func(t *testing.T) {
		ln := &acceptErrorListener{raise: new(temporaryError)}
		s := new(TCPServer)

		closed := make(chan struct{})
		go func() {
			err := s.Serve(ln)
			if g, w := err, ErrServerClosed; g != w {
				t.Errorf("Serve returns got %v, want %v", g, w)
			}
			close(closed)
		}()

		select {
		case <-closed:
			t.Fatal("server got closed, want sleep and call next Accept()")
		case <-time.After(50 * time.Millisecond):
			// ok
		}

		s.Close()

		select {
		case <-closed:
			// test ok
		case <-time.After(3 * time.Second):
			t.Errorf("timeout exceeded while waiting for serv Close")
		}
	})

	t.Run("other error", func(t *testing.T) {
		want := errors.New("other error")
		ln := &acceptErrorListener{raise: want}
		s := new(TCPServer)

		closed := make(chan struct{})
		go func() {
			got := s.Serve(ln)
			if g, w := got.Error(), want.Error(); g != w {
				t.Errorf("Serve returns got %v, want %v", g, w)
			}
			close(closed)
		}()

		select {
		case <-closed:
			// test ok
		case <-time.After(3 * time.Second):
			t.Errorf("timeout exceeded while waiting for serv Close")
		}
	})
}