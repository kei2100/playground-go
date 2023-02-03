package http

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOverwriteTimeoutRoundTripper(t *testing.T) {
	t.Parallel()

	defaultTimeout := time.Second
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: defaultTimeout,
	}
	timeoutPerPath := map[string]time.Duration{
		"/overwrite": time.Millisecond,
	}
	r := OverwriteTimeoutRoundTripper{
		RoundTripper:   &transport,
		TimeoutPerPath: timeoutPerPath,
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(defaultTimeout * 3)
		w.WriteHeader(200)
	}))
	t.Cleanup(svr.Close)

	t.Run("default timeout", func(t *testing.T) {
		t.Parallel()

		path, _ := url.JoinPath(svr.URL, "default")
		req, _ := http.NewRequest("GET", path, nil)
		begin := time.Now()
		_, err := r.RoundTrip(req)
		assert.Error(t, err)
		assert.WithinRange(t, time.Now(), begin.Add(defaultTimeout), begin.Add(defaultTimeout+100*time.Millisecond))
	})
	t.Run("overwrite timeout", func(t *testing.T) {
		t.Parallel()

		path, _ := url.JoinPath(svr.URL, "overwrite")
		req, _ := http.NewRequest("GET", path, nil)
		begin := time.Now()
		_, err := r.RoundTrip(req)
		assert.Error(t, err)
		assert.WithinRange(t, time.Now(), begin.Add(time.Millisecond), begin.Add(time.Millisecond+100*time.Millisecond))
	})
}
