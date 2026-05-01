package auth

import (
	"context"
	"net/http"
)

type ctxKey int

const userKey ctxKey = 1

// UserFromContext returns the User attached to ctx by Required, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

// Middleware bundles the dependencies needed by Required/Optional.
type Middleware struct {
	store      *Store
	cookieName string
}

// NewMiddleware constructs an auth middleware.
func NewMiddleware(store *Store, cookieName string) *Middleware {
	return &Middleware{store: store, cookieName: cookieName}
}

// Required returns a handler that 401s when no valid session cookie is set.
func (m *Middleware) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := m.userFromRequest(r)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, u)))
	})
}

// Optional attaches the user to context if the cookie is valid; otherwise
// passes through unmodified.
func (m *Middleware) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, ok := m.userFromRequest(r); ok {
			r = r.WithContext(context.WithValue(r.Context(), userKey, u))
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) userFromRequest(r *http.Request) (*User, bool) {
	c, err := r.Cookie(m.cookieName)
	if err != nil || c.Value == "" {
		return nil, false
	}
	u, err := m.store.ValidateSession(r.Context(), c.Value)
	if err != nil {
		return nil, false
	}
	return u, true
}
