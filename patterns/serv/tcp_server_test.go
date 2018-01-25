package serv

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"
)

func mustTCPListen(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

func mustTCPDial(t *testing.T, addr net.Addr) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

var nopTCPHandler = func(conn net.Conn) {
	conn.Close()
}

func TestTCPServer_ServeClose(t *testing.T) {
	t.Parallel()

	handler := func(conn net.Conn) {
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
	}

	s := new(TCPServer)
	ln := mustTCPListen(t)

	closed := make(chan struct{})
	go func() {
		err := s.Serve(ln, handler)
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
		t.Fatal("timeout exceeded while waiting for serv Close")
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

func TestTCPServer_Serve_HandleAcceptError(t *testing.T) {
	t.Parallel()

	t.Run("temporary error", func(t *testing.T) {
		ln := &acceptErrorListener{raise: new(temporaryError)}
		s := new(TCPServer)

		closed := make(chan struct{})
		go func() {
			err := s.Serve(ln, nopTCPHandler)
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
			t.Fatal("timeout exceeded while waiting for serv Close")
		}
	})

	t.Run("other error", func(t *testing.T) {
		want := errors.New("other error")
		ln := &acceptErrorListener{raise: want}
		s := new(TCPServer)

		closed := make(chan struct{})
		go func() {
			got := s.Serve(ln, nopTCPHandler)
			if g, w := got.Error(), want.Error(); g != w {
				t.Errorf("Serve returns got %v, want %v", g, w)
			}
			close(closed)
		}()

		select {
		case <-closed:
			// test ok
		case <-time.After(3 * time.Second):
			t.Fatal("timeout exceeded while waiting for serv Close")
		}
	})
}

func TestTCPServer_DoubleServe(t *testing.T) {
	t.Parallel()

	ln := mustTCPListen(t)
	s := new(TCPServer)
	defer s.Close()

	errch := make(chan error, 2)

	// 片方すぐエラーになること
	go func() {
		errch <- s.Serve(ln, nopTCPHandler)
	}()
	go func() {
		errch <- s.Serve(ln, nopTCPHandler)
	}()

	select {
	case err := <-errch:
		if err == nil {
			t.Errorf("got nil, want an error")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout exceeded while waiting for serv Close")
	}
}

func TestTCPServer_DoubleClose(t *testing.T) {
	t.Parallel()

	s := new(TCPServer)
	defer s.Close()

	wg := sync.WaitGroup{}
	closeFunc := func(num int) {
		defer wg.Done()
		if err := s.Close(); err != nil {
			t.Errorf("num %v got %v, want nil", num, err)
		}
	}

	wg.Add(2)
	go closeFunc(1)
	go closeFunc(2)
	wg.Wait()

	ln := mustTCPListen(t)
	done := make(chan struct{})
	go func() {
		s.Serve(ln, nopTCPHandler)
		close(done)
	}()

	wg.Add(2)
	go closeFunc(3)
	go closeFunc(4)
	wg.Wait()

	select {
	case <-done:
		return
	case <-time.After(3 * time.Second):
		t.Fatal("timeout exceeded while waiting for serv Close")
	}
}

func TestTCPServer_Serve_WithOptions(t *testing.T) {
	t.Parallel()

	serveAndDial := func(o TCPServerOptions, f TCPHandleFunc) *TCPServer {
		ln := mustTCPListen(t)
		s := new(TCPServer)
		go s.Serve(ln, f, o)

		if _, err := net.Dial("tcp", ln.Addr().String()); err != nil {
			t.Fatal(err)
		}
		return s
	}

	t.Run("ReadTimeout", func(t *testing.T) {
		errCh := make(chan error, 1)
		s := serveAndDial(WithReadTimeout(100*time.Millisecond), func(conn net.Conn) {
			defer conn.Close()
			_, err := ioutil.ReadAll(conn)
			errCh <- err
		})
		defer s.Close()

		select {
		case err := <-errCh:
			if err == nil {
				t.Error("errCh got nil, want an error")
			}
			// ok
		case <-time.After(1 * time.Second):
			t.Fatal("timeout exceeded while waiting for serv Close")
		}
	})
}

func TestTCPServer_Stats(t *testing.T) {
	t.Parallel()

	t.Run("NumConnections", func(t *testing.T) {
		ln := mustTCPListen(t)
		s := new(TCPServer)
		defer s.Close()

		var serveIn, serveOut sync.WaitGroup
		gate := make(chan struct{})

		go s.Serve(ln, func(conn net.Conn) {
			serveIn.Done()
			<-gate
			conn.Close()
			serveOut.Done()
		})

		const n = 10
		serveIn.Add(n)
		serveOut.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				conn := mustTCPDial(t, ln.Addr())
				<-gate
				conn.Close()
			}()
		}

		testDone := make(chan struct{})
		go func() {
			defer close(testDone)
			serveIn.Wait()

			if g, w := s.Stats().NumConnections, n; g != w {
				t.Errorf("NumConnections got %v, want %v", g, w)
			}

			close(gate)
			serveOut.Wait()
		}()

		select {
		case <-testDone:
			break
		case <-time.After(3 * time.Second):
			t.Fatal("timeout exceeded while waiting for serv Close")
		}
		if g, w := s.Stats().NumConnections, 0; g != w {
			t.Errorf("NumConnections got %v, want %v", g, w)
		}
	})
}
