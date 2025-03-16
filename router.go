// Package rahjoo provides a flexible and lightweight routing solution for the standard library's `http.ServeMux`.
// It allows you to define routes with HTTP methods, group routes under common prefixes, and apply middlewares
// to handlers. The package is designed to work seamlessly with the standard `net/http` library, making it
// easy to integrate into existing projects.
//
// Key Features:
// - Define routes with HTTP methods (e.g., GET, POST).
// - Group routes under common prefixes (e.g., "/api/v1").
// - Apply middlewares to individual handlers or groups of routes.
// - Combine multiple route groups into a single routing table.
// - Bind routes to the standard `http.ServeMux` for compatibility with the `net/http` library.
//
// Example Usage:
//
//	// Define a handler function
//	helloHandler := func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	}
//
//	// Create a GroupRoute with a prefix
//	group := router.GroupRoute{
//	    "/api/v1": router.Route{
//	        "/hello": {
//	            http.MethodGet: router.NewHandler(helloHandler, middleware.Recovery),
//	        },
//	    },
//	}
//
//	// Convert the GroupRoute to a Route
//	routes := router.NewGroup(group)
//
//	// Bind the routes to the ServeMux
//	mux := http.NewServeMux()
//	router.BindRoutesToMux(mux, routes)
//
//	// Start the HTTP server
//	http.ListenAndServe(":8080", mux)
//
// For more details, see the documentation for individual types and functions.
package rahjoo

import (
	"fmt"
	"net/http"

	"github.com/amirzayi/rahjoo/middleware"
)

type (
	// actionHandler is a struct that encapsulates an HTTP handler function along with
	// a list of middlewares to be applied to it. This allows for flexible and reusable
	// middleware chaining for specific routes or actions.
	actionHandler struct {
		// handler is the main HTTP handler function that will process the request.
		// It is the core logic for handling an HTTP request after all middlewares have been applied.
		handler http.HandlerFunc
		// middlewares is a slice of middleware.Middleware functions that will be applied
		// to the handler. These middlewares are executed in the order they are defined,
		// with the last middleware in the slice being the first to execute (closest to the handler).
		middlewares []middleware.Middleware
	}

	// Path represents the URL path for a route (e.g., "/shelves/{shelf_id}/books").
	// It is used to define the endpoint for an HTTP request.
	Path string

	// Method represents the HTTP method for a request (e.g., http.MethodGet, http.MethodPost).
	// If left empty, it will handle all HTTP methods for the given path.
	Method string

	// Route defines a mapping of URL paths to their corresponding HTTP methods and handlers.
	// It is a map where:
	// - The key is a Path (URL path).
	// - The value is another map where:
	// - The key is a Method (HTTP method).
	// - The value is an actionHandler(function to handle the request).
	// This structure allows for flexible route definitions with support for multiple HTTP methods per path.
	Route map[Path]map[Method]actionHandler

	// Prefix represents a common prefix for a group of routes.
	// It is used to group related routes under a shared URL prefix (e.g., "/api/v1").
	Prefix string

	// GroupRoute defines a mapping of route prefixes to their corresponding Route maps.
	// It is a map where:
	// - The key is a Prefix (common URL prefix).
	// - The value is a Route (collection of paths and their handlers).
	// This allows for organizing routes into logical groups based on their prefixes.
	GroupRoute map[Prefix]Route
)

// NewGroup creates a new Route map by combining a GroupRoute with optional middlewares.
// It iterates over the GroupRoute, prepends the prefix to each path, and applies the provided middlewares
// to all handlers in the group. This is useful for grouping routes under a common prefix (e.g., "/api/v1").
func NewGroup(group GroupRoute, middlewares ...middleware.Middleware) Route {
	routes := Route{}
	for prefix, route := range group {
		for path, method := range route {
			routes[Path(fmt.Sprintf("%s%s", prefix, path))] = method
		}
	}
	return routes
}

// NewHandler creates an actionHandler with the given HTTP handler and middlewares.
// The middlewares are applied in reverse order, meaning the last middleware in the list
// will be executed first (closest to the handler).
func NewHandler(handler http.HandlerFunc, middlewares ...middleware.Middleware) actionHandler {
	return actionHandler{
		handler:     handler,
		middlewares: middlewares,
	}
}

// BindRoutesToMux binds the provided routes to a http.ServeMux.
// It iterates over the routes, applies the middlewares to each handler using middleware.Chain,
// and registers the handlers with the ServeMux. The route paths are combined with their HTTP methods
// to create unique route identifiers (e.g., "GET /api/v1/books").
func BindRoutesToMux(mux *http.ServeMux, routes ...Route) {
	mergedRoutes := mergeRoutes(routes...)
	for route, handler := range mergedRoutes {
		for method, action := range handler {
			mux.Handle(fmt.Sprintf("%s %s", method, route), middleware.Chain(action.handler, action.middlewares...))
		}
	}
}

// mergeRoutes combines multiple Route maps into a single Route map.
// It iterates over each Route and merges them into a single map, ensuring that
// routes with the same path and method are not overwritten.
func mergeRoutes(routes ...Route) Route {
	merged := Route{}
	for _, route := range routes {
		for path, method := range route {
			merged[path] = method
		}
	}
	return merged
}
