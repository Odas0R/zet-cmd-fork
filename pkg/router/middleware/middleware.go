package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func emptyMiddleware(next http.Handler) http.Handler {
  return next
}
