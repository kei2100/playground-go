package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Route setups routing for this server
func (s *Server) Route() {
	r := newRouter()
	defer func() { s.router = r }()

	r.Use(maxBody(s.MaxBodyBytes))
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}

type router struct {
	chi.Router
}

func newRouter() *router {
	chir := chi.NewRouter()
	return &router{
		Router: chir,
	}
}
