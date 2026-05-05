package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
			writeJSON(w, http.StatusUnauthorized, loginResponse{Message: "Incorrect username or password"})
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

// userListItem is the response shape for GET /users.
type userListItem struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	IsSuperuser bool   `json:"is_superuser"`
	Enabled     bool   `json:"enabled"`
}

type createUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	IsSuperuser bool   `json:"is_superuser"`
}

type updateUserRequest struct {
	FullName    *string `json:"full_name"`
	Email       *string `json:"email"`
	IsSuperuser *bool   `json:"is_superuser"`
	Enabled     *bool   `json:"enabled"`
}

type adminPasswordResetRequest struct {
	NewPassword string `json:"new_password"`
}

type selfPasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ListUsers returns all users. Requires superuser.
func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.cfg.Store.ListUsers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
		return
	}
	out := make([]userListItem, 0, len(users))
	for _, u := range users {
		out = append(out, userListItem{
			Username:    u.Username,
			FullName:    u.FullName,
			Email:       u.Email,
			IsSuperuser: u.IsSuperuser,
			Enabled:     u.Enabled,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

// CreateUser creates a new user. Requires superuser.
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "invalid body"})
		return
	}
	if req.Username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "username is required"})
		return
	}
	if req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "password is required"})
		return
	}
	err := h.cfg.Store.CreateUser(r.Context(), CreateUserParams{
		Username:    req.Username,
		Password:    req.Password,
		FullName:    req.FullName,
		Email:       req.Email,
		IsSuperuser: req.IsSuperuser,
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "conflict", "message": "username already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"success": true, "username": req.Username})
}

// UpdateUser updates profile/flags for a user. Requires superuser.
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	caller, _ := UserFromContext(r.Context())

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "invalid body"})
		return
	}

	// Safety rails for self-modification.
	if caller != nil && caller.Username == username {
		if req.IsSuperuser != nil && !*req.IsSuperuser {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "cannot remove your own superuser flag"})
			return
		}
		if req.Enabled != nil && !*req.Enabled {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "cannot disable your own account"})
			return
		}
	}

	if req.FullName != nil || req.Email != nil {
		fn := ""
		em := ""
		if req.FullName != nil {
			fn = *req.FullName
		}
		if req.Email != nil {
			em = *req.Email
		}
		// If only one field was provided, load the current value for the other.
		if req.FullName == nil || req.Email == nil {
			u, err := h.cfg.Store.GetUser(r.Context(), username)
			if errors.Is(err, ErrUserNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
				return
			}
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
				return
			}
			if req.FullName == nil {
				fn = u.FullName
			}
			if req.Email == nil {
				em = u.Email
			}
		}
		if err := h.cfg.Store.UpdateProfile(r.Context(), username, fn, em); err != nil {
			if errors.Is(err, ErrUserNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
			return
		}
	}

	if req.IsSuperuser != nil {
		if err := h.cfg.Store.SetSuperuser(r.Context(), username, *req.IsSuperuser); err != nil {
			if errors.Is(err, ErrUserNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
			return
		}
	}

	if req.Enabled != nil {
		if err := h.cfg.Store.SetEnabled(r.Context(), username, *req.Enabled); err != nil {
			if errors.Is(err, ErrUserNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// DeleteUser deletes a user. Requires superuser.
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	caller, _ := UserFromContext(r.Context())
	if caller != nil && caller.Username == username {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "cannot delete yourself"})
		return
	}
	if err := h.cfg.Store.DeleteUser(r.Context(), username); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AdminResetPassword resets a user's password. Requires superuser.
func (h *Handlers) AdminResetPassword(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	var req adminPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "invalid body"})
		return
	}
	if req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "new_password is required"})
		return
	}
	if err := h.cfg.Store.UpdatePassword(r.Context(), username, req.NewPassword); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found", "message": "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// SelfChangePassword allows an authenticated user to change their own password.
func (h *Handlers) SelfChangePassword(w http.ResponseWriter, r *http.Request) {
	caller, ok := UserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized", "message": "not authenticated"})
		return
	}
	var req selfPasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "invalid body"})
		return
	}
	if req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request", "message": "new_password is required"})
		return
	}
	_, err := h.cfg.Store.VerifyPassword(r.Context(), caller.Username, req.CurrentPassword, h.cfg.Lockout)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized", "message": "current password is incorrect"})
		return
	}
	if err := h.cfg.Store.UpdatePassword(r.Context(), caller.Username, req.NewPassword); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal", "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
