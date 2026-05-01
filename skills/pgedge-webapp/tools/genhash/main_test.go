package main

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateHash(t *testing.T) {
	h, err := generateHash("Hunter2HunterTwo!")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(h, "$2") {
		t.Errorf("not bcrypt: %q", h)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte("Hunter2HunterTwo!")); err != nil {
		t.Errorf("hash does not verify: %v", err)
	}
}

func TestGenerateHash_RejectsEmpty(t *testing.T) {
	if _, err := generateHash(""); err == nil {
		t.Fatal("expected error on empty password")
	}
}
