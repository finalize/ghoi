package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestHealthz(t *testing.T) {
	rec := do(t, http.MethodGet, "/healthz")

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, 期待は 200", rec.Code)
	}
	if got := rec.Body.String(); got != "ok\n" {
		t.Errorf("body = %q, 期待は %q", got, "ok\n")
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q", ct)
	}
}

// ServeMux にメソッドを書いてあるので、GET 以外は 405 になる。
func TestHealthzRejectsOtherMethods(t *testing.T) {
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		t.Run(m, func(t *testing.T) {
			if rec := do(t, m, "/healthz"); rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("code = %d, 期待は 405", rec.Code)
			}
		})
	}
}

func TestUnknownPath(t *testing.T) {
	if rec := do(t, http.MethodGet, "/そんなものはない"); rec.Code != http.StatusNotFound {
		t.Errorf("code = %d, 期待は 404", rec.Code)
	}
}
