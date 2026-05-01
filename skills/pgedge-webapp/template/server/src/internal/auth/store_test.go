package auth

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
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
