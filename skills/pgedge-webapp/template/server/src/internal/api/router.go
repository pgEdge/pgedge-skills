// Package api wires HTTP routes for the server.
package api

import (
	"encoding/json"
	"net/http"

	"<MODULE_PATH>/server/src/internal/auth"
)

// Deps bundles router dependencies.
type Deps struct {
	Store        *auth.Store
	Handlers     *auth.Handlers
	Middleware   *auth.Middleware
	LoginLimiter *auth.RateLimiter
}

// NewRouter returns the configured ServeMux.
func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("POST /api/v1/auth/login",
		rateLimit(d.LoginLimiter, http.HandlerFunc(d.Handlers.Login)))
	mux.HandleFunc("POST /api/v1/auth/logout", d.Handlers.Logout)

	mux.Handle("GET /api/v1/user/info",
		d.Middleware.Required(http.HandlerFunc(d.Handlers.UserInfo)))

	mux.Handle("GET /api/v1/hello",
		d.Middleware.Required(http.HandlerFunc(helloHandler)))

	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.UserFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message":  "Hello, " + u.Username,
		"username": u.Username,
	})
}

func rateLimit(rl *auth.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(clientIP(r)) {
			http.Error(w, `{"error":"rate limited"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := indexComma(fwd); i > 0 {
			return fwd[:i]
		}
		return fwd
	}
	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}
	return "unknown"
}

func indexComma(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			return i
		}
	}
	return -1
}
