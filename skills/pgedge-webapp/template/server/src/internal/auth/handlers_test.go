package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestHandlers(t *testing.T) (*Handlers, *Store) {
	s := newTestStore(t)
	h := NewHandlers(HandlersConfig{
		Store:           s,
		CookieName:      "test_session",
		SessionLifetime: time.Hour,
		Lockout:         LockoutConfig{MaxFailedAttempts: 5, LockoutDuration: time.Minute},
		SecureCookie:    false,
	})
	return h, s
}

func TestLogin_Success(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "Hunter2HunterTwo!"})

	body := strings.NewReader(`{"username":"alice","password":"Hunter2HunterTwo!"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		b, _ := io.ReadAll(rec.Body)
		t.Fatalf("code=%d body=%s", rec.Code, b)
	}
	if !strings.Contains(rec.Header().Get("Set-Cookie"), "test_session=") {
		t.Errorf("no cookie set: %q", rec.Header().Get("Set-Cookie"))
	}
}

func TestLogin_BadCreds(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "Hunter2HunterTwo!"})
	rec := httptest.NewRecorder()
	body := strings.NewReader(`{"username":"alice","password":"nope"}`)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	h.Login(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestLogout_ClearsCookie(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "u", Password: "Hunter2HunterTwo!"})
	tok, _ := s.CreateSession(ctx, "u", time.Hour)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: tok})
	h.Logout(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("code = %d", rec.Code)
	}
	setCookie := rec.Header().Get("Set-Cookie")
	if !strings.Contains(setCookie, "test_session=") || !strings.Contains(setCookie, "Expires=Thu, 01 Jan 1970") {
		t.Errorf("cookie not cleared: %q", setCookie)
	}
	if _, err := s.ValidateSession(context.Background(), tok); err == nil {
		t.Errorf("session should be deleted")
	}
}

func TestUserInfo_Authenticated(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "u", Password: "Hunter2HunterTwo!", IsSuperuser: true})
	tok, _ := s.CreateSession(ctx, "u", time.Hour)

	mw := NewMiddleware(s, "test_session")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/user/info", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: tok})
	mw.Required(http.HandlerFunc(h.UserInfo)).ServeHTTP(rec, req)

	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out["authenticated"] != true || out["username"] != "u" || out["is_superuser"] != true {
		t.Errorf("body = %+v", out)
	}
}
