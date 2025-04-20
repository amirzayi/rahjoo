## Introduction

[![Go Reference](https://pkg.go.dev/badge/github.com/amirzayi/rahjoo.svg)](https://pkg.go.dev/github.com/amirzayi/rahjoo)

Rahjoo(**رهـجـو**) meaning Pathfinder in persian, is a lightweight zero dependency HTTP router library for Go, designed to work seamlessly with the standard `net/http` library. This library allows you to define routes with HTTP methods, group routes under common prefixes, and apply middlewares to handlers.

## Features

- **HTTP Method Routing**: Define routes for specific HTTP methods (e.g., GET, POST).
- **Route Grouping**: Group routes under common prefixes (e.g., `/api/v1`).
- **Middleware Support**: Apply middlewares to individual handlers or groups of routes.
- **Compatibility**: Works with the standard `http.ServeMux` for easy integration.
- **Lightweight**: Minimal dependencies and easy to use.

---

## Installation

To install the library, use `go get`:

```bash
go get github.com/amirzayi/rahjoo
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	chim "github.com/go-chi/chi/v5/middleware"

	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
	"github.com/amirzayi/rahjoo/middleware/cors"
)

func h(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func listUsers(w http.ResponseWriter, r *http.Request) {}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	postUsersRoute := rahjoo.Route{
		"/posts/{id}/users": {
			http.MethodGet: rahjoo.NewHandler(listUsers, chim.NoCache, chim.RequestID),
		},
	}

	userV1Gp := rahjoo.NewGroupRoute("/api/v1/users", rahjoo.Route{
		"/list": {
			http.MethodGet: rahjoo.NewHandler(listUsers, middleware.EnforceJSON),
		},
		"/{id}": {
			http.MethodGet: rahjoo.NewHandler(http.NotFound),
		},
	}).SetMiddleware(chim.NoCache, chim.CleanPath)

	// empty method will route all method on given handler
	userV2Gp := rahjoo.NewGroupRoute("/api/v2", rahjoo.Route{
		"/users": {
			"": rahjoo.NewHandler(listUsers, chim.Throttle(6)),
		},
	}).SetMiddleware(middleware.EnforceJSON, chim.NoCache)

	// bind routes to http multiplexer
	rahjoo.BindRoutesToMux(mux, userV1Gp, userV2Gp, postUsersRoute)

	// you can use middleware developed based on std http.HttpHandler
	// such as chi router middlewares
	handler := middleware.Chain(mux,
		cors.CORSHandler(),
		middleware.Recovery(log.Default()),
		chim.Timeout(time.Second),
		chim.Logger,
		chim.RealIP)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
