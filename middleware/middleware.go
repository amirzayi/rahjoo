// Package middleware provides HTTP middlewares.
package middleware

import (
	"log"
	"mime"
	"net/http"
	"slices"
)

// Middleware is a type that represents an HTTP middleware function.
// It takes an http.Handler and returns a new http.Handler that wraps the original.
type Middleware func(http.Handler) http.Handler

// Chain applies a series of middlewares to an http.Handler.
// The middlewares are applied by order, meaning the last middleware in the list
// will be the last to execute when handling an HTTP request.
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	slices.Reverse(middlewares)
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Recovery is a middleware that recovers from panics during HTTP request handling.
// It logs the panic and returns a 500 Internal Server Error response to the client.
// The logger parameter is used to log the panic details.
func Recovery(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Printf("panic recovered: %v\n", rec)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// EnforceJSON is a middleware that ensures the incoming HTTP request has a Content-Type header
// set to "application/json". If the header is missing or invalid, it returns an appropriate
// error response (400 Bad Request or 415 Unsupported Media Type).
func EnforceJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			http.Error(w, "Content-Type header is not set", http.StatusBadRequest)
			return
		}
		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, "Content-Type header is not set", http.StatusBadRequest)
			return
		}
		if mt != "application/json" {
			http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
