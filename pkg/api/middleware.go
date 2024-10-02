package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader assigns status code and header to ResponseWriter of statusResponseWriter object
func (sw *statusResponseWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.ResponseWriter.WriteHeader(statusCode)
}

// Generate (if there is none) and store request_id both in the request context AND X-REQUEST-ID header
func RequestIDLoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var id string
			if r.URL.Query().Has("request_id") {
				id = r.URL.Query().Get("request_id")
			} else if id = r.Header.Get("x-request-id"); id != "" {
				w.Header().Set("x-request-id", id)
			} else {
				id = uuid.NewString()
			}
			ctx := context.WithValue(r.Context(), "request_id", id)
			r = r.WithContext(ctx)
			w.Header().Set("x-request-id", id)
			next.ServeHTTP(w, r)
		})
	}
}

// Log each HTTP request with its execution time and request_id
func LoggerMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := NewStatusResponseWriter(w)
			defer func() {
				log.Printf(
					"method=%s code=%v time=%v request_id=%s url=%s://%s%s",
					r.Method,
					sw.statusCode,
					time.Since(start),
					r.Context().Value("request_id"),
					r.URL.Scheme,
					r.Host,
					r.RequestURI,
				)
			}()
			next.ServeHTTP(sw, r)
		})
	}
}

func URLSchemaMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// we're not going to terminate SSL in this service
			r.URL.Scheme = "http"

			next.ServeHTTP(w, r)
		})
	}
}
