package proxy

import (
	"io"
	"net/http"
)

// Proxy Server
type Proxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Hop-by-hop headers, which are meaningful only for a single
// transport-level connection, and are not stored by caches or
// forwarded by proxies.
//
// https://tools.ietf.org/html/rfc2616#section-13.5.1
var hopByHopHeaders = map[string]struct{}{
	// Header names are canonicalized (see http.Request or http.Response).
	"Connection":          struct{}{},
	"Keep-Alive":          struct{}{},
	"Proxy-Authenticate":  struct{}{},
	"Proxy-Authorization": struct{}{},
	"TE":                struct{}{},
	"Trailers":          struct{}{},
	"Transfer-Encoding": struct{}{},
	"Upgrade":           struct{}{},
}

type proxy struct {
	transport *http.Client
}

// New returns Proxy
func New() Proxy {
	return &proxy{transport: http.DefaultClient}
}

func (x *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cp := new(http.Request)
	copyRequest(r, cp)
	rewriteRequest(cp)

	res, err := x.transport.Do(cp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for h, vv := range res.Header {
		for _, v := range vv {
			w.Header().Set(h, v)
		}
	}
	io.Copy(w, res.Body)
}

func copyRequest(orig *http.Request, cp *http.Request) {
	cp.Method = orig.Method
	cp.URL = orig.URL
	cp.Header = make(http.Header)
	for k, v := range orig.Header {
		if _, ok := hopByHopHeaders[k]; ok {
			continue
		}
		cp.Header[k] = v
	}
	cp.Body = orig.Body
}

func rewriteRequest(r *http.Request) {
	r.URL.Host = "www.google.jp"
	r.URL.Scheme = "https"
}
