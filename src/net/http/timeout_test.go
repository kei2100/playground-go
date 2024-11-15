package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"golang.org/x/net/nettest"

	"github.com/stretchr/testify/assert"
)

type sleepReader struct {
	d time.Duration
	r io.Reader
}

func (r *sleepReader) Read(p []byte) (n int, err error) {
	time.Sleep(r.d)
	return r.r.Read(p)
}

func TestHTTPRequestContextTimeout(t *testing.T) {
	timeout := time.Millisecond * 100

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	t.Cleanup(svr.Close)

	t.Run("will timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// create request
		path, _ := url.JoinPath(svr.URL, "/")
		body := &sleepReader{
			d: timeout * 2,
			r: bytes.NewBufferString(""),
		}
		req, _ := http.NewRequestWithContext(ctx, "GET", path, body)
		// send request
		_, err := http.DefaultClient.Do(req)
		assert.Error(t, err)
		assert.Truef(t, errors.Is(err, context.DeadlineExceeded), "err is %T", err)
	})
	t.Run("will not timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// create request
		path, _ := url.JoinPath(svr.URL, "/ignore_req_body")
		body := &sleepReader{
			d: timeout / 2,
			r: bytes.NewBufferString(""),
		}
		req, _ := http.NewRequestWithContext(ctx, "GET", path, body)
		// send request
		_, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
	})
}

func TestResponseHeaderTimeout(t *testing.T) {
	t.Parallel()
	t.SkipNow()

	t.Run("will timeout", func(t *testing.T) {
		t.Parallel()

		body := &sleepReader{
			d: time.Second,
			r: bytes.NewBufferString("hello"),
		}
		transport := &http.Transport{
			ResponseHeaderTimeout: time.Millisecond * 100,
		}
		cli := &http.Client{Transport: transport}
		req, _ := http.NewRequest("POST", "https://httpbin.org/delay/1", body)
		_, err := cli.Do(req)
		assert.Error(t, err)
	})
	t.Run("will not timeout", func(t *testing.T) {
		t.Parallel()

		body := &sleepReader{
			d: time.Second,
			r: bytes.NewBufferString("hello"),
		}
		transport := &http.Transport{
			ResponseHeaderTimeout: time.Second * 2,
		}
		cli := &http.Client{Transport: transport}
		req, _ := http.NewRequest("POST", "https://httpbin.org/delay/1", body)
		_, err := cli.Do(req)
		assert.NoError(t, err)
	})
}

func TestHTTPServerReadTimeout(t *testing.T) {
	timeout := time.Millisecond * 500
	ln, err := nettest.NewLocalListener("tcp")
	if err != nil {
		t.Fatal(err)
	}
	svr := http.Server{
		ReadTimeout: timeout,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Timeout() {
					w.WriteHeader(408)
					return
				}
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
		}),
	}
	t.Cleanup(func() {
		svr.Close()
	})
	go svr.Serve(ln)

	t.Run("will timeout", func(t *testing.T) {
		body := &sleepReader{
			d: timeout * 2,
			r: bytes.NewBufferString("hello"),
		}
		req, _ := http.NewRequest("POST", "http://"+ln.Addr().String(), body)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 408, resp.StatusCode)
	})
	t.Run("will not timeout", func(t *testing.T) {
		body := &sleepReader{
			d: timeout / 50,
			r: bytes.NewBufferString("hello"),
		}
		req, _ := http.NewRequest("POST", "http://"+ln.Addr().String(), body)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})
}
