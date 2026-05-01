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

	mw := NewMiddleware(s, "test_session")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: tok})

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
	mw := NewMiddleware(s, "test_session")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	mw.Required(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestOptional_PassesThroughWithNoUser(t *testing.T) {
	s := newTestStore(t)
	mw := NewMiddleware(s, "test_session")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	called := false
	mw.Optional(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		_, ok := UserFromContext(r.Context())
		if ok {
			t.Error("no user expected in context")
		}
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
	if !called {
		t.Error("handler not called")
	}
}

func TestOptional_AttachesUserWhenValid(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "opt", Password: "Hunter2HunterTwo!"})
	tok, _ := s.CreateSession(ctx, "opt", time.Hour)

	mw := NewMiddleware(s, "test_session")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "test_session", Value: tok})

	mw.Optional(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext(r.Context())
		if !ok || u.Username != "opt" {
			t.Errorf("expected opt user, got ok=%v u=%+v", ok, u)
		}
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
}
