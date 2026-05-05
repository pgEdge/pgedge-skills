package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPath_UserConfigDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	userPath := filepath.Join(tmp, "pgedge", "<BINARY_NAME>.yaml")
	if err := os.MkdirAll(filepath.Dir(userPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userPath, []byte("http: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, source := GetDefaultConfigPath("<BINARY_NAME>.yaml")
	if got != userPath {
		t.Errorf("path = %q, want %q", got, userPath)
	}
	if source != "user" {
		t.Errorf("source = %q, want %q", source, "user")
	}
}

func TestGetDefaultConfigPath_FallsThroughToSystem(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	got, source := GetDefaultConfigPath("<BINARY_NAME>.yaml")
	want := filepath.Join("/etc/pgedge", "<BINARY_NAME>.yaml")
	if got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if source != "system" {
		t.Errorf("source = %q, want %q", source, "system")
	}
}

func TestFileExists(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "x")
	if FileExists(p) {
		t.Fatal("missing file reports as existing")
	}
	if err := os.WriteFile(p, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	if !FileExists(p) {
		t.Fatal("present file reports as missing")
	}
}
