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
