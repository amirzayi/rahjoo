package cors_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/amirzayi/rahjoo/middleware/cors"
)

const routeCorsPath = "/api/cors"

func TestCors(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
		w.WriteHeader(http.StatusOK)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(routeCorsPath, h)

	corshandler := cors.CORSHandler()
	handler := corshandler(mux)

	for _, tc := range []struct {
		name,
		method,
		expectedBody string
		expectedStatus  int
		headers         map[string]string
		expectedHeaders map[string][]string
	}{
		{"Success", http.MethodGet, "hello", http.StatusOK, nil, nil},
		{"Cors Success", http.MethodOptions, "", http.StatusNoContent, map[string]string{
			"Origin":                        "*",
			"Access-Control-Request-Method": http.MethodGet,
		}, map[string][]string{
			"Access-Control-Allow-Origin":  cors.DefaultOriginAllowList,
			"Access-Control-Allow-Methods": cors.DefaultMethodAllowList,
			"Access-Control-Allow-Headers": cors.DefaultHeadersAllowList,
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, routeCorsPath, http.NoBody)
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Result().StatusCode != tc.expectedStatus {
				t.Fatalf("expected %d, got %d", tc.expectedStatus, rec.Result().StatusCode)
			}

			if body := rec.Body.String(); body != tc.expectedBody {
				t.Fatalf("expected %s, got %s", tc.expectedBody, body)
			}

			for k, v := range tc.expectedHeaders {
				headers := strings.Split(rec.Header().Get(k), ", ")

				if len(headers) != len(v) {
					t.Fatalf("expected response %d headers, got %d", len(v), len(headers))
				}

				if slices.Compare(headers, v) != 0 {
					t.Fatalf("expected response headers %v, got %v", v, headers)
				}
			}
		})
	}
}
