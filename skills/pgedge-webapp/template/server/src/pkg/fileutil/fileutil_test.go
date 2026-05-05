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

// TestGetDefaultConfigPath_UserConfigDirFallback exercises the branch where
// XDG_CONFIG_HOME is set but has no matching file, and the file exists under
// os.UserConfigDir() instead.
func TestGetDefaultConfigPath_UserConfigDirFallback(t *testing.T) {
	// Point XDG_CONFIG_HOME at an empty temp dir so that branch is entered
	// but no file is found there.
	xdgTmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdgTmp)

	// Discover the real UserConfigDir and plant the file there.
	userDir, err := os.UserConfigDir()
	if err != nil {
		t.Skip("os.UserConfigDir() unavailable:", err)
	}
	targetDir := filepath.Join(userDir, "pgedge")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configFile := filepath.Join(targetDir, "<BINARY_NAME>-fallback-test.yaml")
	if err := os.WriteFile(configFile, []byte("http: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(configFile) })

	got, source := GetDefaultConfigPath("<BINARY_NAME>-fallback-test.yaml")
	if got != configFile {
		t.Errorf("path = %q, want %q", got, configFile)
	}
	if source != SourceUser {
		t.Errorf("source = %q, want %q", source, SourceUser)
	}
}
