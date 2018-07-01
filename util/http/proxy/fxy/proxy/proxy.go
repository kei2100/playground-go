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
	Transport   *http.Client
	once        sync.Once
	Destination *url.URL
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.once.Do(func() {
		if s.Transport == nil {
			s.Transport = http.DefaultClient
		}
	})

	cp := new(http.Request)
	s.copyRequest(r, cp)

	res, err := s.Transport.Do(cp)
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

func (s *Server) copyRequest(orig, cp *http.Request) {
	cp.Method = orig.Method

	cp.URL = new(url.URL)
	s.copyURL(orig.URL, cp.URL)

	cp.Header = make(http.Header)
	for k, v := range orig.Header {
		if _, ok := hopByHopHeaders[k]; ok {
			continue
		}
		cp.Header[k] = v
	}

	cp.Body = orig.Body
}

func (s *Server) copyURL(orig, cp *url.URL) {
	if s.Destination == nil {
		panic("proxy: Server.Destination must be set")
	}

	cp.Scheme = s.Destination.Scheme
	cp.Host = s.Destination.Host

	if s.Destination.User != nil {
		cp.User = s.Destination.User
	} else {
		cp.User = orig.User
	}

	if dplen, oplen := len(s.Destination.Path), len(orig.Path); dplen > 0 && oplen > 0 {
		cp.Path = s.Destination.Path + "/" + orig.Path
	} else if dplen > 0 && oplen == 0 {
		cp.Path = s.Destination.Path
	} else if dplen == 0 && oplen > 0 {
		cp.Path = orig.Path
	}

	if dplen, oplen := len(s.Destination.RawPath), len(orig.RawPath); dplen > 0 && oplen > 0 {
		cp.RawPath = s.Destination.RawPath + "/" + orig.RawPath
	} else if dplen > 0 && oplen == 0 {
		cp.RawPath = s.Destination.RawPath
	} else if dplen == 0 && oplen > 0 {
		cp.RawPath = orig.RawPath
	}

	cp.ForceQuery = orig.ForceQuery
	cp.RawQuery = orig.RawQuery
	cp.Fragment = orig.Fragment
}
