package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	something = "something"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		ctx := context.WithValue(r.Context(), something, time.Since(start))
		req := r.WithContext(ctx)

		next.ServeHTTP(wrapped, req)

		log.Println(r.Method, r.URL.Path, time.Since(start), wrapped.statusCode)
	})
}
