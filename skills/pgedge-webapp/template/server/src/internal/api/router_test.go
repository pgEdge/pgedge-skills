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

// ---- clientIP tests ----

func TestClientIP_XForwardedForSingle(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	if got := clientIP(req); got != "1.2.3.4" {
		t.Errorf("clientIP = %q", got)
	}
}

func TestClientIP_XForwardedForMultiple(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	if got := clientIP(req); got != "1.2.3.4" {
		t.Errorf("clientIP = %q", got)
	}
}

func TestClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "9.8.7.6:1234"
	if got := clientIP(req); got != "9.8.7.6:1234" {
		t.Errorf("clientIP = %q", got)
	}
}

func TestClientIP_Unknown(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = ""
	if got := clientIP(req); got != "unknown" {
		t.Errorf("clientIP = %q", got)
	}
}

// ---- indexComma tests ----

func TestIndexComma_NoComma(t *testing.T) {
	if got := indexComma("abc"); got != -1 {
		t.Errorf("indexComma = %d", got)
	}
}

func TestIndexComma_LeadingComma(t *testing.T) {
	if got := indexComma(",abc"); got != 0 {
		t.Errorf("indexComma = %d", got)
	}
}

func TestIndexComma_Mid(t *testing.T) {
	if got := indexComma("abc,def"); got != 3 {
		t.Errorf("indexComma = %d", got)
	}
}

// ---- rateLimit middleware test ----

func TestRateLimit_Blocks(t *testing.T) {
	deps := testDeps(t)
	ctx := context.Background()
	_ = deps.Store.CreateUser(ctx, auth.CreateUserParams{Username: "u2", Password: "Hunter2HunterTwo!"})

	// Use a limiter that allows only 1 request per window.
	strictLimiter := auth.NewRateLimiter(1, time.Minute)
	deps.LoginLimiter = strictLimiter

	r := NewRouter(deps)

	doLogin := func() int {
		body := strings.NewReader(`{"username":"u2","password":"Hunter2HunterTwo!"}`)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		return rec.Code
	}

	// First request should succeed.
	if code := doLogin(); code != http.StatusOK {
		t.Fatalf("first login code = %d", code)
	}
	// Second request should be rate-limited.
	if code := doLogin(); code != http.StatusTooManyRequests {
		t.Fatalf("second login code = %d, want 429", code)
	}
}

// ---- UserInfo without auth ----

func TestRouter_UserInfoWithoutAuth(t *testing.T) {
	r := NewRouter(testDeps(t))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("GET", "/api/v1/user/info", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- Logout without cookie ----

func TestRouter_LogoutNoCookie(t *testing.T) {
	r := NewRouter(testDeps(t))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest("POST", "/api/v1/auth/logout", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- Login with malformed JSON ----

func TestRouter_LoginMalformedJSON(t *testing.T) {
	r := NewRouter(testDeps(t))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}
