package http

import (
	"net/http"

	"github.com/kei2100/playground-go/patterns/serv/http/internal/response"
)

func maxBody(bytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > bytes {
				response.SendPayloadTooLarge(w)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
