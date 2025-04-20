package cors_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amirzayi/rahjoo/middleware/cors"
)

func ExampleCORSHandler() {
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
		w.WriteHeader(http.StatusOK)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/test", h)

	corshandler := cors.CORSHandler(
		cors.WithHeaders(append(cors.DefaultHeadersAllowList, "Api-Key", "CSRF-Token")),
		cors.WithMethods(append(cors.DefaultMethodAllowList, http.MethodHead)),
		cors.WithOrigins([]string{"http://127.0.0.1:9090", "http://127.0.0.1:9091"}),
	)
	handler := corshandler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func ExampleWithMethods() {
	cors.WithMethods([]string{http.MethodGet, http.MethodPost})
}

func ExampleWithHeaders() {
	cors.WithHeaders([]string{"Origin", "Authorization"})
}

func ExampleWithOrigins() {
	cors.WithOrigins([]string{"*"})
}
