package main

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"<MODULE_PATH>/server/src/internal/auth"
)

func newStore(t *testing.T) *auth.Store {
	t.Helper()
	s, err := auth.Open(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestAddUser_RequiresUsername(t *testing.T) {
	s := newStore(t)
	f, _ := ParseFlags([]string{"-add-user"}, "")
	if err := runUserCommand(context.Background(), s, f, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestAddUser_OK(t *testing.T) {
	s := newStore(t)
	f, _ := ParseFlags([]string{
		"-add-user", "-username", "alice",
		"-password", "Hunter2HunterTwo!",
		"-email", "alice@example.com",
	}, "")
	if err := runUserCommand(context.Background(), s, f, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, err := s.GetUser(context.Background(), "alice")
	if err != nil || u.Email != "alice@example.com" {
		t.Errorf("u=%+v err=%v", u, err)
	}
}

func TestListUsers_Output(t *testing.T) {
	s := newStore(t)
	_ = s.CreateUser(context.Background(), auth.CreateUserParams{
		Username: "alice", Password: "Hunter2HunterTwo!",
	})
	f, _ := ParseFlags([]string{"-list-users"}, "")
	var buf bytes.Buffer
	if err := runUserCommand(context.Background(), s, f, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "alice") {
		t.Errorf("output: %q", buf.String())
	}
}

func TestDeleteUser_OK(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, auth.CreateUserParams{Username: "del", Password: "Hunter2HunterTwo!"})
	f, _ := ParseFlags([]string{"-delete-user", "-username", "del"}, "")
	if err := runUserCommand(ctx, s, f, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
}

func TestEnableDisableUser(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, auth.CreateUserParams{Username: "eu", Password: "Hunter2HunterTwo!"})

	f, _ := ParseFlags([]string{"-disable-user", "-username", "eu"}, "")
	if err := runUserCommand(ctx, s, f, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, _ := s.GetUser(ctx, "eu")
	if u.Enabled {
		t.Error("expected disabled")
	}

	f2, _ := ParseFlags([]string{"-enable-user", "-username", "eu"}, "")
	if err := runUserCommand(ctx, s, f2, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, _ = s.GetUser(ctx, "eu")
	if !u.Enabled {
		t.Error("expected enabled")
	}
}

func TestSetUnsetSuperuser(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, auth.CreateUserParams{Username: "su", Password: "Hunter2HunterTwo!"})

	f, _ := ParseFlags([]string{"-set-superuser", "-username", "su"}, "")
	if err := runUserCommand(ctx, s, f, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, _ := s.GetUser(ctx, "su")
	if !u.IsSuperuser {
		t.Error("expected superuser")
	}

	f2, _ := ParseFlags([]string{"-unset-superuser", "-username", "su"}, "")
	if err := runUserCommand(ctx, s, f2, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, _ = s.GetUser(ctx, "su")
	if u.IsSuperuser {
		t.Error("expected non-superuser")
	}
}

func TestUpdateUser_ProfileAndPassword(t *testing.T) {
	s := newStore(t)
	ctx := context.Background()
	_ = s.CreateUser(ctx, auth.CreateUserParams{Username: "up", Password: "Hunter2HunterTwo!"})

	// Update full name and email
	f, _ := ParseFlags([]string{"-update-user", "-username", "up", "-full-name", "Up User", "-email", "up@example.com"}, "")
	if err := runUserCommand(ctx, s, f, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	u, _ := s.GetUser(ctx, "up")
	if u.FullName != "Up User" || u.Email != "up@example.com" {
		t.Errorf("u = %+v", u)
	}

	// Update password
	f2, _ := ParseFlags([]string{"-update-user", "-username", "up", "-password", "NewPass99!"}, "")
	if err := runUserCommand(ctx, s, f2, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	// Verify new password works
	uu, err := s.VerifyPassword(ctx, "up", "NewPass99!", auth.LockoutConfig{MaxFailedAttempts: 5})
	if err != nil || uu.Username != "up" {
		t.Errorf("new password verify: %v", err)
	}
}

func TestUpdateUser_RequiresUsername(t *testing.T) {
	s := newStore(t)
	f, _ := ParseFlags([]string{"-update-user"}, "")
	if err := runUserCommand(context.Background(), s, f, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunUserCommand_NoCommand(t *testing.T) {
	s := newStore(t)
	f := &Flags{} // no fs, no user command set
	f2, _ := ParseFlags(nil, "")
	_ = f
	if err := runUserCommand(context.Background(), s, f2, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error for no command")
	}
}

func TestIsSet(t *testing.T) {
	f, _ := ParseFlags([]string{"-addr", ":8888", "-debug"}, "")
	if !f.IsSet("addr") {
		t.Error("addr should be set")
	}
	if !f.IsSet("debug") {
		t.Error("debug should be set")
	}
	if f.IsSet("tls") {
		t.Error("tls should not be set")
	}
	if f.IsSet("cert") {
		t.Error("cert should not be set")
	}
}
