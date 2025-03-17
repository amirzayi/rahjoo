package rahjoo_test

import (
	"log"
	"net/http"

	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
)

func ExampleNewGroup() {
	listUsers := func(http.ResponseWriter, *http.Request) {}

	userV1 := rahjoo.GroupRoute{"/api/v1/users": {
		"/list": {
			http.MethodGet: rahjoo.NewHandler(listUsers, middleware.EnforceJSON),
		},
		"/{id}": {
			http.MethodGet: rahjoo.NewHandler(http.NotFound),
		},
	}}

	_ = rahjoo.NewGroup(userV1, middleware.Recovery(log.Default()))
}

func ExampleNewHandler() {
	h := func(w http.ResponseWriter, r *http.Request) {}

	rahjoo.NewHandler(h, middleware.EnforceJSON, middleware.Recovery(log.Default()))
}

func ExampleBindRoutesToMux() {
	listUsers := func(http.ResponseWriter, *http.Request) {}

	userV1 := rahjoo.GroupRoute{"/api/v1/users": {
		"/list": {
			http.MethodGet: rahjoo.NewHandler(listUsers, middleware.EnforceJSON),
		},
		"/{id}": {
			http.MethodGet: rahjoo.NewHandler(http.NotFound),
		},
	}}
	userV1Gp := rahjoo.NewGroup(userV1, middleware.Recovery(log.Default()))

	userV2Gp := rahjoo.NewGroup(rahjoo.GroupRoute{"/api/v2": {
		"/users": {
			"": rahjoo.NewHandler(listUsers),
		},
	}}, middleware.EnforceJSON, middleware.Recovery(log.Default()))

	mux := http.NewServeMux()

	rahjoo.BindRoutesToMux(mux, userV1Gp, userV2Gp)
}
