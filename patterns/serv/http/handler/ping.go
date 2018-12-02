package handler

import (
	"net/http"

	"github.com/kei2100/playground-go/patterns/serv/http/internal/response"
)

// Ping handler
func Ping() http.HandlerFunc {
	type res struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		response.SendJSON(w, 200, &res{Message: "ok"})
	}
}
