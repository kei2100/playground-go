package go1_20

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
)

func TestHTTPResponseController_ReadDeadline(t *testing.T) {
	t.Parallel()

	const (
		defaultReadTimeout = time.Millisecond * 500
		shortenReadTimeout = time.Millisecond
		extendReadTimeout  = time.Second * 2
	)
	// create mux
	mux := http.NewServeMux()
	mux.HandleFunc("/read/default", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
	mux.HandleFunc("/read/shorten", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)
		rc.SetReadDeadline(time.Now().Add(shortenReadTimeout))
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
	mux.HandleFunc("/read/extend", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)
		rc.SetReadDeadline(time.Now().Add(extendReadTimeout))
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
	// create server
	ln, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	svr := http.Server{
		Handler:     mux,
		ReadTimeout: defaultReadTimeout,
	}
	go svr.Serve(ln)
	t.Cleanup(func() { svr.Close() })
	// tests
	t.Run("default, ok", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Post(fmt.Sprintf("http://%s/read/default", ln.Addr()), "text/plain", &waitedReader{wait: 100 * time.Millisecond})
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("default, timeout", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Post(fmt.Sprintf("http://%s/read/default", ln.Addr()), "text/plain", &waitedReader{wait: time.Second})
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})
	t.Run("shorten, timeout", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Post(fmt.Sprintf("http://%s/read/shorten", ln.Addr()), "text/plain", &waitedReader{wait: 100 * time.Millisecond})
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})
	t.Run("extend, ok", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Post(fmt.Sprintf("http://%s/read/extend", ln.Addr()), "text/plain", &waitedReader{wait: time.Second})
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestHTTPResponseController_WriteDeadline(t *testing.T) {
	t.Parallel()

	const (
		defaultWriteTimeout = time.Millisecond * 500
		shortenWriteTimeout = time.Millisecond
		extendWriteTimeout  = time.Second * 2
	)
	// create mux
	mux := http.NewServeMux()
	mux.HandleFunc("/write/default/ok", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, &waitedReader{wait: 100 * time.Millisecond})
	})
	mux.HandleFunc("/write/default/timeout", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, &waitedReader{wait: time.Second})
	})
	mux.HandleFunc("/write/shorten/timeout", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)
		rc.SetWriteDeadline(time.Now().Add(shortenWriteTimeout))
		io.Copy(w, &waitedReader{wait: 100 * time.Millisecond})
	})
	mux.HandleFunc("/write/extend/ok", func(w http.ResponseWriter, r *http.Request) {
		rc := http.NewResponseController(w)
		rc.SetWriteDeadline(time.Now().Add(extendWriteTimeout))
		io.Copy(w, &waitedReader{wait: time.Second})
	})
	// create server
	ln, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	svr := http.Server{
		Handler:      mux,
		WriteTimeout: defaultWriteTimeout,
	}
	go svr.Serve(ln)
	t.Cleanup(func() { svr.Close() })
	// tests
	t.Run("default, ok", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Get(fmt.Sprintf("http://%s/write/default/ok", ln.Addr()))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("default, timeout", func(t *testing.T) {
		t.Parallel()
		_, err := http.Get(fmt.Sprintf("http://%s/write/default/timeout", ln.Addr()))
		assert.Truef(t, errors.Is(err, io.EOF), "err is %v", err) // EOF が返却される
	})
	t.Run("shorten, timeout", func(t *testing.T) {
		t.Parallel()
		_, err := http.Get(fmt.Sprintf("http://%s/write/shorten/timeout", ln.Addr()))
		assert.Truef(t, errors.Is(err, io.EOF), "err is %v", err)
	})
	t.Run("extend, ok", func(t *testing.T) {
		t.Parallel()
		resp, err := http.Get(fmt.Sprintf("http://%s/write/extend/ok", ln.Addr()))
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

type waitedReader struct {
	wait time.Duration
}

func (w waitedReader) Read(_ []byte) (n int, err error) {
	time.Sleep(w.wait)
	return 0, io.EOF
}
