package middleware

import (
	"log"
	"net/http"
	"time"
)

// wrappedResponseWriter adds statusCode to ResponseWriter
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader intercepts the call to write the header to capture the status code.
func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// WithLogger logs each request with the response status and duration.
func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}
