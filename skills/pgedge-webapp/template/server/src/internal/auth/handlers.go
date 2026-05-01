package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// HandlersConfig collects the dependencies of NewHandlers.
type HandlersConfig struct {
	Store           *Store
	CookieName      string
	SessionLifetime time.Duration
	Lockout         LockoutConfig
	SecureCookie    bool
}

// Handlers exposes the auth HTTP handlers.
type Handlers struct{ cfg HandlersConfig }

// NewHandlers constructs the handler set.
func NewHandlers(cfg HandlersConfig) *Handlers { return &Handlers{cfg: cfg} }

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginResponse struct {
	Success   bool   `json:"success"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Login authenticates a username/password and sets the session cookie.
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, loginResponse{Message: "invalid body"})
		return
	}
	u, err := h.cfg.Store.VerifyPassword(r.Context(), req.Username, req.Password, h.cfg.Lockout)
	if err != nil {
		switch {
		case errors.Is(err, ErrAccountLocked):
			writeJSON(w, http.StatusTooManyRequests, loginResponse{Message: "account locked"})
		case errors.Is(err, ErrAccountDisabled):
			writeJSON(w, http.StatusForbidden, loginResponse{Message: "account disabled"})
		default:
			writeJSON(w, http.StatusUnauthorized, loginResponse{Message: "invalid credentials"})
		}
		return
	}
	token, err := h.cfg.Store.CreateSession(r.Context(), u.Username, h.cfg.SessionLifetime)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, loginResponse{Message: "session error"})
		return
	}
	expires := time.Now().Add(h.cfg.SessionLifetime)
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   h.cfg.SecureCookie,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, loginResponse{
		Success:   true,
		ExpiresAt: expires.UTC().Format(time.RFC3339),
	})
}

// Logout clears the session cookie and deletes the matching session row.
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(h.cfg.CookieName); err == nil && c.Value != "" {
		_ = h.cfg.Store.DeleteSession(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   h.cfg.SecureCookie,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// UserInfo returns the authenticated user; requires Required middleware.
func (h *Handlers) UserInfo(w http.ResponseWriter, r *http.Request) {
	u, ok := UserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"authenticated": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated": true,
		"username":      u.Username,
		"is_superuser":  u.IsSuperuser,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
