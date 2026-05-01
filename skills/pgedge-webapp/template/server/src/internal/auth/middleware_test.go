package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequired_AllowsValidCookie(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "u", Password: "Hunter2HunterTwo!"})
	tok, _ := s.CreateSession(ctx, "u", time.Hour)

	mw := NewMiddleware(s, "<COOKIE_NAME>")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "<COOKIE_NAME>", Value: tok})

	handler := mw.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := UserFromContext(r.Context())
		if u == nil || u.Username != "u" {
			t.Fatalf("user not in context: %+v", u)
		}
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestRequired_RejectsMissingCookie(t *testing.T) {
	s := newTestStore(t)
	mw := NewMiddleware(s, "<COOKIE_NAME>")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	mw.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}
