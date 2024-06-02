package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/odas0r/zet/pkg/view"
)

type ResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

func WithLayout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is an htmx request
		if r.Header.Get("HX-Request") == "true" {
			next.ServeHTTP(w, r)
			return
		}

		crw := &ResponseWriter{
			ResponseWriter: w,
			buf:            new(bytes.Buffer),
		}

		next.ServeHTTP(crw, r)

		html := crw.buf.String()

		component := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := w.Write([]byte(html))
			return err
		})
		layoutComponent := view.Layout(component)
		templ.Handler(layoutComponent).ServeHTTP(w, r)
	})
}
