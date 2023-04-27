package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPathBasedTimeoutRoundTripper(t *testing.T) {
	t.Parallel()

	timeout := 100 * time.Millisecond
	rt := &PathBasedTimeoutRoundTripper{
		PathTimeout: map[string]time.Duration{
			"/foo": timeout / 2,
		},
	}

	t.Run("path based timeout", func(t *testing.T) {
		t.Parallel()
		// create server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(timeout)
			w.WriteHeader(200)
		}))
		t.Cleanup(svr.Close)
		// send request
		client := &http.Client{
			Transport: rt,
		}
		_, err := client.Get(svr.URL + "/foo")
		assert.Truef(t, errors.Is(err, context.DeadlineExceeded), "err is %T", err)
	})
	t.Run("no timeout", func(t *testing.T) {
		t.Parallel()
		// create server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(timeout)
			w.WriteHeader(200)
		}))
		t.Cleanup(svr.Close)
		// send request
		client := &http.Client{
			Transport: rt,
		}
		resp, err := client.Get(svr.URL + "/bar")
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("default timeout", func(t *testing.T) {
		t.Parallel()
		rt := &PathBasedTimeoutRoundTripper{
			PathTimeout: map[string]time.Duration{
				"/foo": timeout / 2,
			},
			DefaultTimeout: timeout / 2,
		}
		// create server
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(timeout)
			w.WriteHeader(200)
		}))
		t.Cleanup(svr.Close)
		// send request
		client := &http.Client{
			Transport: rt,
		}
		_, err := client.Get(svr.URL + "/bar")
		assert.Truef(t, errors.Is(err, context.DeadlineExceeded), "err is %T", err)
	})
}
