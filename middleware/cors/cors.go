package cors

import (
	"net/http"
	"slices"
	"strings"
)

var DefaultOriginAllowList = []string{"*"}

var DefaultMethodAllowList = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

var DefaultHeadersAllowList = []string{
	"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization", "Origin",
}

type corsHandler struct {
	allowedMethods,
	allowedOrigins,
	allowedHeaders []string
}

func newCorsHandler() *corsHandler {
	return &corsHandler{
		allowedMethods: DefaultMethodAllowList,
		allowedOrigins: DefaultOriginAllowList,
		allowedHeaders: DefaultHeadersAllowList,
	}
}

func isPreflight(r *http.Request) bool {
	return r.Method == http.MethodOptions &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

func (c *corsHandler) hasOrigin(origin string) bool {
	return slices.Contains(c.allowedOrigins, "*") || slices.Contains(c.allowedOrigins, origin)
}

func (c *corsHandler) hasMethod(method string) bool {
	return slices.Contains(c.allowedMethods, method) || method == http.MethodOptions
}

func (c *corsHandler) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		w.Header().Add("Vary", "Origin")

		if isPreflight(r) {
			requestedMethod := r.Header.Get("Access-Control-Request-Method")
			if c.hasMethod(requestedMethod) && c.hasOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.allowedHeaders, ", "))
				w.WriteHeader(http.StatusNoContent)
			}
			return
		}

		if c.hasOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		next.ServeHTTP(w, r)
	})
}

func CORSHandler(opts ...optionCorsFunc) func(next http.Handler) http.Handler {
	cors := newCorsHandler()
	for _, opt := range opts {
		opt(cors)
	}
	return cors.handler
}

type optionCorsFunc func(*corsHandler)

func WithMethods(methods []string) optionCorsFunc {
	return func(ch *corsHandler) {
		ch.allowedMethods = methods
	}
}

func WithHeaders(headers []string) optionCorsFunc {
	return func(ch *corsHandler) {
		ch.allowedHeaders = headers
	}
}

func WithOrigins(origins []string) optionCorsFunc {
	return func(ch *corsHandler) {
		ch.allowedOrigins = origins
	}
}
