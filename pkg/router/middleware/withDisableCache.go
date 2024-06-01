package middleware

import (
	"net/http"
)

func WithDisableCache(enable bool) Middleware {
	if !enable {
		return emptyMiddleware
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			next.ServeHTTP(w, r)
		})
	}
}
