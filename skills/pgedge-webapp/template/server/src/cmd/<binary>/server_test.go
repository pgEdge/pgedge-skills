package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunServer_LogsResolvedConfigSource exercises the source label by
// running runServer in -init mode (which exits cleanly without binding a
// port) and asserting the structured log contains the expected source.
func TestRunServer_LogsResolvedConfigSource(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "cfg.yaml")
	if err := os.WriteFile(cfgPath, []byte("data_dir: "+tmp+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	pwFile := filepath.Join(tmp, "admin.pw")
	if err := os.WriteFile(pwFile, []byte("Hunter2HunterTwo!"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := runServer(context.Background(), []string{
		"-config", cfgPath,
		"-admin-password-file", pwFile,
		"-init",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runServer: %v\nstderr=%s", err, stderr.String())
	}
	if !strings.Contains(stderr.String(), `"source":"cli"`) {
		t.Errorf("source label missing: %s", stderr.String())
	}
}

// TestRunServer_ListUsers exercises running runServer with -list-users after -init.
func TestRunServer_ListUsers(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "cfg.yaml")
	if err := os.WriteFile(cfgPath, []byte("data_dir: "+tmp+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	pwFile := filepath.Join(tmp, "admin.pw")
	if err := os.WriteFile(pwFile, []byte("Hunter2HunterTwo!"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Init first
	var stdout, stderr bytes.Buffer
	if err := runServer(context.Background(), []string{
		"-config", cfgPath,
		"-admin-password-file", pwFile,
		"-init",
	}, &stdout, &stderr); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Now list-users
	stdout.Reset()
	stderr.Reset()
	if err := runServer(context.Background(), []string{
		"-config", cfgPath,
		"-list-users",
	}, &stdout, &stderr); err != nil {
		t.Fatalf("list-users: %v", err)
	}
	if !strings.Contains(stdout.String(), "admin") {
		t.Errorf("admin not in output: %s", stdout.String())
	}
}

// TestRunServer_DefaultConfigPath exercises the default config source when no config file exists.
func TestRunServer_DefaultConfigPath(t *testing.T) {
	tmp := t.TempDir()
	pwFile := filepath.Join(tmp, "admin.pw")
	if err := os.WriteFile(pwFile, []byte("Hunter2HunterTwo!"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Use data-dir to avoid the default path being in the exe dir
	var stdout, stderr bytes.Buffer
	err := runServer(context.Background(), []string{
		"-data-dir", tmp,
		"-admin-password-file", pwFile,
		"-init",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("runServer: %v\nstderr=%s", err, stderr.String())
	}
	// Source should be "defaults" since no config file exists
	if !strings.Contains(stderr.String(), `"source":"defaults"`) {
		t.Errorf("source label missing: %s", stderr.String())
	}
}

// TestRunServer_BadFlags exercises an error from ParseFlags.
func TestRunServer_BadFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := runServer(context.Background(), []string{"-unknown-flag-xyz"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

// TestRunServer_BadConfig exercises an error from LoadConfig.
func TestRunServer_BadConfig(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "bad.yaml")
	if err := os.WriteFile(cfgPath, []byte(":\tbad yaml\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	err := runServer(context.Background(), []string{
		"-config", cfgPath,
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for bad config")
	}
}
