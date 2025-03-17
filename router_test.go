package rahjoo_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amirzayi/rahjoo"
	"github.com/amirzayi/rahjoo/middleware"
)

func TestRouting(t *testing.T) {
	const routePathTest = "/api/test"
	const routePathV2 = "/api/v2"

	h := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case http.MethodPut, http.MethodPatch:
			w.WriteHeader(http.StatusAccepted)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}

	r := rahjoo.NewGroup(rahjoo.GroupRoute{
		routePathTest: {
			"": {
				"": rahjoo.NewHandler(h),
			},
			"/sample": {
				http.MethodPut: rahjoo.NewHandler(h),
			},
			"/sample/{id}/simple": {
				http.MethodPost: rahjoo.NewHandler(h),
			},
			"/not_found": {
				http.MethodGet:    rahjoo.NewHandler(http.NotFound),
				http.MethodDelete: rahjoo.NewHandler(h),
			},
		},
	})

	r2 := rahjoo.NewGroup(rahjoo.GroupRoute{
		routePathV2: {
			"": {
				"": rahjoo.NewHandler(h),
			},
		},
	})

	mux := http.NewServeMux()
	rahjoo.BindRoutesToMux(mux, r, r2)

	testCases := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, routePathTest, http.StatusOK},
		{http.MethodPost, routePathTest, http.StatusCreated},
		{http.MethodPut, routePathTest, http.StatusAccepted},
		{http.MethodDelete, routePathTest, http.StatusNoContent},
		{http.MethodOptions, routePathTest, http.StatusOK},
		{http.MethodGet, fmt.Sprintf("%s%s", routePathTest, "/sample"), http.StatusMethodNotAllowed},
		{http.MethodPut, fmt.Sprintf("%s%s", routePathTest, "/sample"), http.StatusAccepted},
		{http.MethodGet, fmt.Sprintf("%s%s", routePathTest, "/sample/1/simple"), http.StatusMethodNotAllowed},
		{http.MethodPost, fmt.Sprintf("%s%s", routePathTest, "/sample/1/simple"), http.StatusCreated},
		{http.MethodGet, fmt.Sprintf("%s%s", routePathTest, "/not_found"), http.StatusNotFound},
		{http.MethodDelete, fmt.Sprintf("%s%s", routePathTest, "/not_found"), http.StatusNoContent},
		{http.MethodGet, fmt.Sprintf("%s%s", routePathTest, "/not_exists"), http.StatusNotFound},
		{http.MethodGet, routePathV2, http.StatusOK},
		{http.MethodPost, routePathV2, http.StatusCreated},
		{http.MethodPut, routePathV2, http.StatusAccepted},
		{http.MethodDelete, routePathV2, http.StatusNoContent},
		{http.MethodOptions, routePathV2, http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("path %q, method %q", tc.path, tc.method), func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.status {
				t.Errorf("got status code %d, want %d", res.StatusCode, tc.status)
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	panicHandler := func(w http.ResponseWriter, r *http.Request) {
		panic("hi")
	}
	jsonHandler := func(w http.ResponseWriter, r *http.Request) {}

	panicRoute := rahjoo.Route{"/panic": {http.MethodGet: rahjoo.NewHandler(panicHandler, middleware.Recovery(log.Default()))}}
	jsonRoute := rahjoo.Route{"/json": {http.MethodGet: rahjoo.NewHandler(jsonHandler, middleware.EnforceJSON)}}

	mux := http.NewServeMux()
	rahjoo.BindRoutesToMux(mux, panicRoute, jsonRoute)

	testCases := []struct {
		name   string
		path   string
		header string
		status int
	}{
		{"must_recover_panic", "/panic", "application/json", http.StatusInternalServerError},
		{"enforce_json", "/json", "", http.StatusBadRequest},
		{"pass_json", "/json", "application/json", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.path, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tc.header)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.status {
				t.Errorf("got status code %d, want %d", res.StatusCode, tc.status)
			}
		})
	}
}
