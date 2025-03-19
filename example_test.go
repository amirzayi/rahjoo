package rahjoo_test

import (
	"log"
	"net/http"

	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
)

func ExampleNewGroupRoute() {
	listUsers := func(http.ResponseWriter, *http.Request) {}

	_ = rahjoo.NewGroupRoute("/api/v1", rahjoo.Route{
		"/users": {
			http.MethodGet: rahjoo.NewHandler(listUsers),
		},
		"/post/{id}/users": {
			http.MethodGet: rahjoo.NewHandler(listUsers, middleware.Recovery(log.Default())),
		},
	}).SetMiddleware(middleware.EnforceJSON)
}

func ExampleNewHandler() {
	h := func(w http.ResponseWriter, r *http.Request) {}

	rahjoo.NewHandler(h, middleware.EnforceJSON, middleware.Recovery(log.Default()))
}

func ExampleBindRoutesToMux() {
	listUsers := func(http.ResponseWriter, *http.Request) {}

	userV1Gp := rahjoo.NewGroupRoute("/api/v1/users", rahjoo.Route{
		"/list": {
			http.MethodGet: rahjoo.NewHandler(listUsers, middleware.EnforceJSON),
		},
		"/{id}": {
			http.MethodGet: rahjoo.NewHandler(http.NotFound),
		},
	})

	userV2Gp := rahjoo.NewGroupRoute("/api/v2", rahjoo.Route{
		"/users": {
			"": rahjoo.NewHandler(listUsers),
		},
	}).SetMiddleware(middleware.EnforceJSON, middleware.Recovery(log.Default()))

	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	rahjoo.BindRoutesToMux(mux, userV1Gp, userV2Gp)
}
