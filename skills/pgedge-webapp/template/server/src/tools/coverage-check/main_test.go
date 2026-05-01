package main

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleProfile = `mode: atomic
example.com/foo/bar.go:1.1,2.1 1 1
example.com/foo/bar.go:3.1,4.1 1 0
example.com/cmd/<BINARY_NAME>/main.go:1.1,2.1 1 0
`

func TestCoverage_Threshold(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "c.out")
	if err := os.WriteFile(p, []byte(sampleProfile), 0o600); err != nil {
		t.Fatal(err)
	}
	pct, err := computeCoverage(p, []string{"cmd/<BINARY_NAME>/main.go"})
	if err != nil {
		t.Fatal(err)
	}
	if pct != 50.0 {
		t.Errorf("pct = %v, want 50", pct)
	}
}
