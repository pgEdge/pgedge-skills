package auth

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "auth.db")
	s, err := Open(p)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestOpen_CreatesSchema(t *testing.T) {
	s := newTestStore(t)
	n, err := s.CountUsers(context.Background())
	if err != nil {
		t.Fatalf("CountUsers: %v", err)
	}
	if n != 0 {
		t.Errorf("fresh db: count = %d, want 0", n)
	}
}

func TestStore_CreateAndGetUser(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	err := s.CreateUser(ctx, CreateUserParams{
		Username: "alice", Password: "Hunter2HunterTwo!",
		FullName: "Alice Example", Email: "alice@example.com",
		IsSuperuser: true,
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	u, err := s.GetUser(ctx, "alice")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if u.Username != "alice" || !u.IsSuperuser || !u.Enabled {
		t.Errorf("got = %+v", u)
	}
}

func TestStore_DuplicateRejected(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	p := CreateUserParams{Username: "bob", Password: "Hunter2HunterTwo!"}
	if err := s.CreateUser(ctx, p); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateUser(ctx, p); err == nil {
		t.Fatal("duplicate accepted")
	}
}

func TestStore_ListUsers(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "a", Password: "Hunter2HunterTwo!"})
	_ = s.CreateUser(ctx, CreateUserParams{Username: "b", Password: "Hunter2HunterTwo!"})
	list, err := s.ListUsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("len = %d", len(list))
	}
}

func TestStore_DeleteUser(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "c", Password: "Hunter2HunterTwo!"})
	if err := s.DeleteUser(ctx, "c"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetUser(ctx, "c"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestStore_EnableDisable(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "d", Password: "Hunter2HunterTwo!"})
	if err := s.SetEnabled(ctx, "d", false); err != nil {
		t.Fatal(err)
	}
	u, _ := s.GetUser(ctx, "d")
	if u.Enabled {
		t.Error("disable did not stick")
	}
}

func TestVerify_OK(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "carol", Password: "Right-Password-99"})
	cfg := LockoutConfig{MaxFailedAttempts: 3, LockoutDuration: time.Minute}
	u, err := s.VerifyPassword(ctx, "carol", "Right-Password-99", cfg)
	if err != nil || u.Username != "carol" {
		t.Fatalf("verify: u=%+v err=%v", u, err)
	}
}

func TestVerify_Wrong(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "carol", Password: "Right-Password-99"})
	_, err := s.VerifyPassword(ctx, "carol", "wrong", LockoutConfig{MaxFailedAttempts: 3})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v", err)
	}
}

func TestVerify_LocksAfterMax(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "carol", Password: "Right-Password-99"})
	cfg := LockoutConfig{MaxFailedAttempts: 2, LockoutDuration: time.Minute}
	_, _ = s.VerifyPassword(ctx, "carol", "x", cfg)
	_, _ = s.VerifyPassword(ctx, "carol", "x", cfg)
	_, err := s.VerifyPassword(ctx, "carol", "Right-Password-99", cfg)
	if !errors.Is(err, ErrAccountLocked) {
		t.Fatalf("err = %v", err)
	}
}

func TestVerify_DisabledRejected(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "x", Password: "Right-Password-99"})
	_ = s.SetEnabled(ctx, "x", false)
	_, err := s.VerifyPassword(ctx, "x", "Right-Password-99", LockoutConfig{MaxFailedAttempts: 5})
	if !errors.Is(err, ErrAccountDisabled) {
		t.Fatalf("err = %v", err)
	}
}

func TestSessions(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "alice", Password: "Hunter2HunterTwo!"})

	token, err := s.CreateSession(ctx, "alice", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("empty token")
	}
	u, err := s.ValidateSession(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "alice" {
		t.Errorf("got = %+v", u)
	}
	if err := s.DeleteSession(ctx, token); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ValidateSession(ctx, token); !errors.Is(err, ErrSessionInvalid) {
		t.Errorf("err = %v", err)
	}
}

func TestSessions_Expired(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "x", Password: "Hunter2HunterTwo!"})
	tok, _ := s.CreateSession(ctx, "x", -time.Hour)
	if _, err := s.ValidateSession(ctx, tok); !errors.Is(err, ErrSessionInvalid) {
		t.Errorf("err = %v", err)
	}
	n, err := s.DeleteExpiredSessions(ctx)
	if err != nil || n != 1 {
		t.Errorf("delete expired: n=%d err=%v", n, err)
	}
}

func TestStore_CreateUserWithHash(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	// Use any syntactically valid bcrypt hash string; we only test DB storage,
	// not bcrypt verification here.
	hash := "$2a$12$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012345"
	if err := s.CreateUserWithHash(ctx, "hashuser", hash, true); err != nil {
		t.Fatalf("CreateUserWithHash: %v", err)
	}
	// Verify the user is in the DB by counting.
	n, err := s.CountUsers(ctx)
	if err != nil {
		t.Fatalf("CountUsers: %v", err)
	}
	if n != 1 {
		t.Errorf("count = %d, want 1", n)
	}
	// Duplicate should fail.
	if err := s.CreateUserWithHash(ctx, "hashuser", hash, false); err == nil {
		t.Error("duplicate should be rejected")
	}
}

func TestStore_UpdateProfile(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "p", Password: "Hunter2HunterTwo!"})
	if err := s.UpdateProfile(ctx, "p", "Full Name", "p@example.com"); err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	u, _ := s.GetUser(ctx, "p")
	if u.FullName != "Full Name" || u.Email != "p@example.com" {
		t.Errorf("u = %+v", u)
	}
	// Non-existent user should return ErrUserNotFound
	if err := s.UpdateProfile(ctx, "nobody", "x", "x"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestStore_UpdatePassword(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "pw", Password: "Hunter2HunterTwo!"})
	if err := s.UpdatePassword(ctx, "pw", "NewPass99!"); err != nil {
		t.Fatalf("UpdatePassword: %v", err)
	}
	// New password should work
	u, err := s.VerifyPassword(ctx, "pw", "NewPass99!", LockoutConfig{MaxFailedAttempts: 5})
	if err != nil || u.Username != "pw" {
		t.Errorf("verify after update: u=%+v err=%v", u, err)
	}
	// Old password should fail
	_, err = s.VerifyPassword(ctx, "pw", "Hunter2HunterTwo!", LockoutConfig{MaxFailedAttempts: 5})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("old password err = %v", err)
	}
}

func TestStore_SetSuperuser(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, CreateUserParams{Username: "su", Password: "Hunter2HunterTwo!"})
	if err := s.SetSuperuser(ctx, "su", true); err != nil {
		t.Fatalf("SetSuperuser: %v", err)
	}
	u, _ := s.GetUser(ctx, "su")
	if !u.IsSuperuser {
		t.Error("expected superuser")
	}
	if err := s.SetSuperuser(ctx, "su", false); err != nil {
		t.Fatal(err)
	}
	u, _ = s.GetUser(ctx, "su")
	if u.IsSuperuser {
		t.Error("expected non-superuser")
	}
}

func TestStore_DeleteUser_NotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	if err := s.DeleteUser(ctx, "nobody"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestValidateSession_Empty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	if _, err := s.ValidateSession(ctx, ""); !errors.Is(err, ErrSessionInvalid) {
		t.Errorf("err = %v", err)
	}
}

func TestVerify_UnknownUser(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_, err := s.VerifyPassword(ctx, "nobody", "pass", LockoutConfig{MaxFailedAttempts: 3})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("err = %v", err)
	}
}
