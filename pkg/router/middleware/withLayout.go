package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/odas0r/zet/pkg/view"
)

// WithLayout is a middleware that adds the layout to the component
// if the component is a ComponentHandler
func WithLayout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
