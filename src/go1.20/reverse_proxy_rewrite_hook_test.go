package go1_20

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseProxyRewriteHook(t *testing.T) {
	// target
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	t.Cleanup(targetServer.Close)
	targetURL, err := url.Parse(targetServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	// proxy
	proxy := &httputil.ReverseProxy{
		// Use ReverseProxy.Rewrite instead of Director
		Director: nil,
		Rewrite: func(r *httputil.ProxyRequest) {
			// add Host header
			r.SetURL(targetURL)
			r.Out.Host = targetURL.Host
			// add XFF header (also X-Forwarded-Host, X-Forwarded-Proto)
			//
			// Note: If the outbound request contains an existing X-Forwarded-For header,
			// SetXForwarded appends the client IP address to it. To append to the
			// inbound request's X-Forwarded-For header (the default behavior of
			// ReverseProxy when using a Director function), copy the header
			// from the inbound request before calling SetXForwarded:
			//	rewriteFunc := func(r *httputil.ProxyRequest) {
			//		r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
			//		r.SetXForwarded()
			//	}
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
			r.SetXForwarded()
		},
	}
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))
	t.Cleanup(proxyServer.Close)
	// send request to the proxy
	req, _ := http.NewRequest("GET", proxyServer.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
}
