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
