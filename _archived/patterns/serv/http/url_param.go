package http

import (
	"net/http"

	"github.com/go-chi/chi"
)

type (
	// ParamHandlerFunc is a HandlerFunc that accepts one URL parameter
	ParamHandlerFunc func(w http.ResponseWriter, r *http.Request, param string)
	// Param2HandlerFunc is a HandlerFunc that accepts two URL parameters
	Param2HandlerFunc func(w http.ResponseWriter, r *http.Request, param1, param2 string)
)

func withURLParam(handler ParamHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		if rctx == nil {
			handler(w, r, "")
			return
		}
		vv := rctx.URLParams.Values
		if len(vv) == 0 {
			handler(w, r, "")
			return
		}
		handler(w, r, vv[0])
	}
}

func withURLParam2(handler Param2HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		if rctx == nil {
			handler(w, r, "", "")
			return
		}
		vv := rctx.URLParams.Values
		var v1, v2 string
		if len(vv) > 0 {
			v1 = vv[0]
		}
		if len(vv) > 1 {
			v2 = vv[1]
		}
		handler(w, r, v1, v2)
	}
}
