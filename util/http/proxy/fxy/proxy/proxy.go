package proxy

import (
	"io"
	"net/http"
	"net/url"
	"sync"
)

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

// Server is a forward proxy server
type Server struct {
	Config    Config
	Forwarder *http.Client
	once      sync.Once
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.once.Do(func() {
		if s.Forwarder == nil {
			// TODO keep-alive
			t := &http.Transport{TLSClientConfig: s.Config.TLSClientConfig.TLSConfig()}
			s.Forwarder = &http.Client{Transport: t}
		}
	})

	cp := new(http.Request)
	s.copyRequest(r, cp)
	s.rewriteHeader(cp)
	s.rewriteURL(cp.URL)

	res, err := s.Forwarder.Do(cp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for h, vv := range res.Header {
		for _, v := range vv {
			w.Header().Add(h, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func (s *Server) copyRequest(orig, cp *http.Request) {
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

func (s *Server) rewriteHeader(req *http.Request) {
	for k := range s.Config.Header {
		if k == "Host" {
			req.Host = s.Config.Header.Get(k)
			continue
		}
		req.Header.Set(k, s.Config.Header.Get(k))
	}
}

func (s *Server) rewriteURL(u *url.URL) {
	if len(s.Config.Server) == 0 {
		panic("proxy: Server.Config.Server must be set")
	}
	u.Scheme = s.Config.Scheme()
	u.Host = s.Config.Host()
	u.User = s.Config.UserInfo()
	for _, rewrite := range s.Config.PathRewriters() {
		if ok := rewrite.Do(u); ok {
			break
		}
	}
}
