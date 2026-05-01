package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePassword_FromFlag(t *testing.T) {
	got, err := ResolvePassword("hunter2", true, "")
	if err != nil || got != "hunter2" {
		t.Errorf("got=%q err=%v", got, err)
	}
}

func TestResolvePassword_FromFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "pw")
	if err := os.WriteFile(p, []byte("secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ResolvePassword("", false, p)
	if err != nil || got != "secret" {
		t.Errorf("got=%q err=%v", got, err)
	}
}

func TestResolvePassword_Empty(t *testing.T) {
	got, err := ResolvePassword("", false, "")
	if err != nil || got != "" {
		t.Errorf("got=%q err=%v", got, err)
	}
}
