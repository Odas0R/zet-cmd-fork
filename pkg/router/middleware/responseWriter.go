package middleware

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (crw *ResponseWriter) Write(b []byte) (int, error) {
	return crw.buf.Write(b)
}

func (crw *ResponseWriter) WriteHeader(statusCode int) {
	crw.ResponseWriter.WriteHeader(statusCode)
}

func (crw *ResponseWriter) Header() http.Header {
	return crw.ResponseWriter.Header()
}
