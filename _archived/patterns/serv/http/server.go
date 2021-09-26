package http

import "net/http"

// Server is a http server
type Server struct {
	MaxBodyBytes int64
	router       *router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
