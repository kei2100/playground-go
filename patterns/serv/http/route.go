package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Routes setups routing for this server
func (s *Server) Routes() {
	r := newRouter()
	defer func() { s.router = r }()

	//r.Use(s.maxBody())
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
