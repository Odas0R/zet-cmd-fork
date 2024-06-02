package middleware

import (
	"net/http"
	"strings"
)

// AllowedMethods is a middleware that restricts the HTTP methods to a set of allowed methods.
func WithMethods(allowedMethods ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedMethods))
	for _, method := range allowedMethods {
		allowed[strings.ToUpper(method)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowed[r.Method] {
				w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
