package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/odas0r/zet/pkg/router/middleware"
)

type Route struct {
	Method  string
	Pattern string
}

type Router struct {
	mux        *http.ServeMux
	middleware []middleware.Middleware
	prefix     string
	routes     []Route
}

func New() *Router {
	return &Router{
		mux:        http.NewServeMux(),
		routes:     []Route{},
		middleware: []middleware.Middleware{},
	}
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc, routeMiddleware ...middleware.Middleware) {
	parts := strings.Split(pattern, " ")

	method := parts[0]
	pattern = parts[1]

	pattern = fmt.Sprintf("%s %s%s", method, r.prefix, pattern)
	h := r.applyMiddleware(http.HandlerFunc(handler), routeMiddleware...)
	r.mux.HandleFunc(pattern, h.ServeHTTP)
	r.routes = append(r.routes, Route{Method: method, Pattern: pattern})
}

func (r *Router) Handle(pattern string, handler http.Handler, routeMiddleware ...middleware.Middleware) {
	parts := strings.Split(pattern, " ")

	method := parts[0]
	pattern = parts[1]

	pattern = fmt.Sprintf("%s %s%s", method, r.prefix, pattern)
	r.mux.Handle(pattern, r.applyMiddleware(handler, routeMiddleware...))
	r.routes = append(r.routes, Route{Method: method, Pattern: pattern})
}

func (r *Router) Group(prefix string) *Router {
	return &Router{
		mux:        r.mux,
		middleware: r.middleware,
		prefix:     r.prefix + strings.TrimSuffix(prefix, "/"),
		routes:     r.routes,
	}
}

func (r *Router) Use(middleware ...middleware.Middleware) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) applyMiddleware(handler http.Handler, routeMiddleware ...middleware.Middleware) http.Handler {
	// Apply route-specific middleware first
	for i := len(routeMiddleware) - 1; i >= 0; i-- {
		handler = routeMiddleware[i](handler)
	}
	// Then apply global middleware
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}

func (r *Router) PrintRoutes() {
	for _, route := range r.routes {
		println(route.Pattern)
	}
}
