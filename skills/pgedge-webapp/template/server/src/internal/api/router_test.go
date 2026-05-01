package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"<MODULE_PATH>/server/src/internal/auth"
)

func TestRouter_Health(t *testing.T) {
	r := NewRouter(testDeps(t))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestRouter_HelloRequiresAuth(t *testing.T) {
	r := NewRouter(testDeps(t))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/hello", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestRouter_HelloAuthenticated(t *testing.T) {
	deps := testDeps(t)
	ctx := context.Background()
	_ = deps.Store.CreateUser(ctx, auth.CreateUserParams{Username: "u", Password: "Hunter2HunterTwo!"})
	tok, _ := deps.Store.CreateSession(ctx, "u", time.Hour)

	r := NewRouter(deps)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/hello", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: tok})
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"u"`) {
		t.Errorf("body = %s", rec.Body.String())
	}
}

func testDeps(t *testing.T) Deps {
	t.Helper()
	s, err := auth.Open(t.TempDir() + "/auth.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	h := auth.NewHandlers(auth.HandlersConfig{
		Store: s, CookieName: "test_session",
		SessionLifetime: time.Hour,
	})
	return Deps{
		Store:        s,
		Handlers:     h,
		Middleware:   auth.NewMiddleware(s, "test_session"),
		LoginLimiter: auth.NewRateLimiter(100, time.Minute),
	}
}
