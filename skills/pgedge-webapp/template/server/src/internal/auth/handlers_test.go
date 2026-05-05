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

func TestUserInfo_NotAuthenticated(t *testing.T) {
	h, _ := newTestHandlers(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/user/info", nil)
	// Call UserInfo directly without middleware — no user in context.
	h.UserInfo(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out["authenticated"] != false {
		t.Errorf("body = %+v", out)
	}
}

func TestLogin_MalformedJSON(t *testing.T) {
	h, _ := newTestHandlers(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	h.Login(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestLogin_AccountLocked(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "locked", Password: "Hunter2HunterTwo!"})
	// Exhaust the lockout limit (5 failures configured in newTestHandlers)
	for i := 0; i < 5; i++ {
		body := strings.NewReader(`{"username":"locked","password":"wrong"}`)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")
		h.Login(rec, req)
	}
	// Now the account should be locked
	body := strings.NewReader(`{"username":"locked","password":"Hunter2HunterTwo!"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	h.Login(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestLogin_AccountDisabled(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "disabled", Password: "Hunter2HunterTwo!"})
	_ = s.SetEnabled(ctx, "disabled", false)
	body := strings.NewReader(`{"username":"disabled","password":"Hunter2HunterTwo!"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	h.Login(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestLogout_NoCookie(t *testing.T) {
	h, _ := newTestHandlers(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	// No cookie set
	h.Logout(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- helpers for user-management handler tests ----

// withUser puts a user record directly into a request context (bypasses cookie lookup).
func withUser(r *http.Request, u *User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userKey, u))
}

// mustUser retrieves or creates a user and loads the full record.
func mustUser(t *testing.T, s *Store, username, password string, super bool) *User {
	t.Helper()
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: username, Password: password, IsSuperuser: super})
	u, err := s.GetUser(ctx, username)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	return u
}

// routedReq creates a request whose path values are parsed by a real http.ServeMux.
func routedReq(method, path, pattern string, body string, ctx context.Context) *http.Request {
	var bodyStr *strings.Reader
	if body != "" {
		bodyStr = strings.NewReader(body)
	} else {
		bodyStr = strings.NewReader("")
	}
	inner := httptest.NewRequest(method, path, bodyStr)
	inner = inner.WithContext(ctx)

	var captured *http.Request
	mux := http.NewServeMux()
	mux.HandleFunc(pattern, func(_ http.ResponseWriter, r *http.Request) {
		captured = r
	})
	mux.ServeHTTP(httptest.NewRecorder(), inner)
	if captured == nil {
		return inner
	}
	return captured.WithContext(ctx)
}

// ---- ListUsers ----

func TestListUsers_Happy(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "bob", "Hunter2HunterTwo!", false)

	req := withUser(httptest.NewRequest("GET", "/api/v1/users", nil), admin)
	rec := httptest.NewRecorder()
	h.ListUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
	var out []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("len = %d", len(out))
	}
}

// ---- CreateUser ----

func TestCreateUser_Happy(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	req := routedReq("POST", "/api/v1/users", "POST /api/v1/users",
		`{"username":"newuser","password":"NewPass123!"}`,
		context.WithValue(context.Background(), userKey, admin))
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out["username"] != "newuser" {
		t.Errorf("response = %+v", out)
	}
}

func TestCreateUser_MissingUsername(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	req := withUser(httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(`{"password":"NewPass123!"}`)), admin)
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestCreateUser_MissingPassword(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	req := withUser(httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(`{"username":"bob"}`)), admin)
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestCreateUser_MalformedJSON(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	req := withUser(httptest.NewRequest("POST", "/api/v1/users", strings.NewReader("{bad")), admin)
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	h, s := newTestHandlers(t)
	ctx := context.Background()
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	_ = s.CreateUser(ctx, CreateUserParams{Username: "existing", Password: "Hunter2HunterTwo!"})

	req := withUser(httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(`{"username":"existing","password":"Hunter2HunterTwo!"}`)), admin)
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusConflict {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- UpdateUser ----

func TestUpdateUser_ProfileOnly(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/target", "PATCH /api/v1/users/{username}",
		`{"full_name":"Target User","email":"t@example.com"}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_SingleField_FullNameOnly(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/target", "PATCH /api/v1/users/{username}",
		`{"full_name":"Only Name"}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_SingleField_EmailOnly(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/target", "PATCH /api/v1/users/{username}",
		`{"email":"only@example.com"}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_SuperuserFlag(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/target", "PATCH /api/v1/users/{username}",
		`{"is_superuser":true}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_EnabledFlag(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/target", "PATCH /api/v1/users/{username}",
		`{"enabled":false}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_MalformedJSON(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/admin", "PATCH /api/v1/users/{username}",
		`{bad`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/nobody", "PATCH /api/v1/users/{username}",
		`{"full_name":"Nobody"}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestUpdateUser_NotFound_SuperuserField(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/nobody", "PATCH /api/v1/users/{username}",
		`{"is_superuser":true}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestUpdateUser_NotFound_EnabledField(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/nobody", "PATCH /api/v1/users/{username}",
		`{"enabled":false}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestUpdateUser_SelfRemoveSuperuser(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/admin", "PATCH /api/v1/users/{username}",
		`{"is_superuser":false}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestUpdateUser_SelfDisable(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("PATCH", "/api/v1/users/admin", "PATCH /api/v1/users/{username}",
		`{"enabled":false}`, ctx)
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- DeleteUser ----

func TestDeleteUser_Happy(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "victim", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("DELETE", "/api/v1/users/victim", "DELETE /api/v1/users/{username}", "", ctx)
	rec := httptest.NewRecorder()
	h.DeleteUser(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("DELETE", "/api/v1/users/nobody", "DELETE /api/v1/users/{username}", "", ctx)
	rec := httptest.NewRecorder()
	h.DeleteUser(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestDeleteUser_SelfDelete(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("DELETE", "/api/v1/users/admin", "DELETE /api/v1/users/{username}", "", ctx)
	rec := httptest.NewRecorder()
	h.DeleteUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- AdminResetPassword ----

func TestAdminResetPassword_Happy(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)
	mustUser(t, s, "target", "Hunter2HunterTwo!", false)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("POST", "/api/v1/users/target/password",
		"POST /api/v1/users/{username}/password",
		`{"new_password":"BrandNew123!"}`, ctx)
	rec := httptest.NewRecorder()
	h.AdminResetPassword(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestAdminResetPassword_EmptyPassword(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("POST", "/api/v1/users/target/password",
		"POST /api/v1/users/{username}/password",
		`{"new_password":""}`, ctx)
	rec := httptest.NewRecorder()
	h.AdminResetPassword(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestAdminResetPassword_MalformedJSON(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("POST", "/api/v1/users/target/password",
		"POST /api/v1/users/{username}/password",
		`{bad`, ctx)
	rec := httptest.NewRecorder()
	h.AdminResetPassword(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestAdminResetPassword_NotFound(t *testing.T) {
	h, s := newTestHandlers(t)
	admin := mustUser(t, s, "admin", "Hunter2HunterTwo!", true)

	ctx := context.WithValue(context.Background(), userKey, admin)
	req := routedReq("POST", "/api/v1/users/nobody/password",
		"POST /api/v1/users/{username}/password",
		`{"new_password":"BrandNew123!"}`, ctx)
	rec := httptest.NewRecorder()
	h.AdminResetPassword(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("code = %d", rec.Code)
	}
}

// ---- SelfChangePassword ----

func TestSelfChangePassword_Happy(t *testing.T) {
	h, s := newTestHandlers(t)
	alice := mustUser(t, s, "alice", "OldPass123!", false)

	req := withUser(httptest.NewRequest("POST", "/api/v1/user/password",
		strings.NewReader(`{"current_password":"OldPass123!","new_password":"NewPass456!"}`)), alice)
	rec := httptest.NewRecorder()
	h.SelfChangePassword(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d body = %s", rec.Code, rec.Body.String())
	}
}

func TestSelfChangePassword_WrongCurrent(t *testing.T) {
	h, s := newTestHandlers(t)
	alice := mustUser(t, s, "alice", "OldPass123!", false)

	req := withUser(httptest.NewRequest("POST", "/api/v1/user/password",
		strings.NewReader(`{"current_password":"WrongPass!","new_password":"NewPass456!"}`)), alice)
	rec := httptest.NewRecorder()
	h.SelfChangePassword(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestSelfChangePassword_EmptyNew(t *testing.T) {
	h, s := newTestHandlers(t)
	alice := mustUser(t, s, "alice", "OldPass123!", false)

	req := withUser(httptest.NewRequest("POST", "/api/v1/user/password",
		strings.NewReader(`{"current_password":"OldPass123!","new_password":""}`)), alice)
	rec := httptest.NewRecorder()
	h.SelfChangePassword(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestSelfChangePassword_MalformedJSON(t *testing.T) {
	h, s := newTestHandlers(t)
	alice := mustUser(t, s, "alice", "OldPass123!", false)

	req := withUser(httptest.NewRequest("POST", "/api/v1/user/password",
		strings.NewReader(`{bad`)), alice)
	rec := httptest.NewRecorder()
	h.SelfChangePassword(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("code = %d", rec.Code)
	}
}

func TestSelfChangePassword_NoAuth(t *testing.T) {
	h, _ := newTestHandlers(t)
	body := strings.NewReader(`{"current_password":"OldPass123!","new_password":"NewPass456!"}`)
	req := httptest.NewRequest("POST", "/api/v1/user/password", body)
	rec := httptest.NewRecorder()
	h.SelfChangePassword(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("code = %d", rec.Code)
	}
}
