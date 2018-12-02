package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kei2100/playground-go/patterns/serv/http/handler"
)

// Route setups routing for this server
func (s *Server) Route() {
	r := newRouter()
	defer func() { s.router = r }()

	r.Use(maxBody(s.MaxBodyBytes))
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Get("/ping", handler.Ping())
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
