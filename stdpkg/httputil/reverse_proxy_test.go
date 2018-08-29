package httputil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	u, _ := url.Parse(upstream.URL)
	rp := httputil.NewSingleHostReverseProxy(u)
	proxy := httptest.NewServer(rp)
	defer proxy.Close()

	res, err := http.Get(proxy.URL)
	if err != nil {
		t.Fatalf("http request failure: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("http read response body failure: %v", err)
	}
	if g, w := string(body), "ok"; g != w {
		t.Errorf("body got %v, want %v", g, w)
	}
}
