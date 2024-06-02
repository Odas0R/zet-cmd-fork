package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/odas0r/zet/pkg/view"
)

// preRender renders the component and returns the output as a new templ.Component.
func preRender(ctx context.Context, component templ.Component) (templ.Component, error) {
	var buffer bytes.Buffer
	err := component.Render(ctx, &buffer)
	if err != nil {
		return nil, err
	}
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write(buffer.Bytes())
		return err
	}), nil
}

// WithLayout is a middleware that adds the layout to the component
// if the component is a ComponentHandler
func WithLayout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pageHandler, ok := next.(*templ.ComponentHandler); ok {
			ctx := r.Context()
			renderedComponent, err := preRender(ctx, pageHandler.Component)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			componentWithLayout := view.Layout(renderedComponent)
			templ.Handler(componentWithLayout).ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
