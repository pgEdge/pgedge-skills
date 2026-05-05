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

func TestCoverage_NoExcludes_AllCovered(t *testing.T) {
	profile := "mode: atomic\nexample.com/pkg/a.go:1.1,2.1 3 1\n"
	tmp := t.TempDir()
	p := filepath.Join(tmp, "c.out")
	if err := os.WriteFile(p, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}
	pct, err := computeCoverage(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if pct != 100.0 {
		t.Errorf("pct = %v, want 100", pct)
	}
}

func TestCoverage_EmptyProfile(t *testing.T) {
	profile := "mode: atomic\n"
	tmp := t.TempDir()
	p := filepath.Join(tmp, "c.out")
	if err := os.WriteFile(p, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}
	pct, err := computeCoverage(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	// total==0 → 100%
	if pct != 100.0 {
		t.Errorf("pct = %v, want 100", pct)
	}
}

func TestCoverage_MissingFile(t *testing.T) {
	_, err := computeCoverage("/no/such/file.out", nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCoverage_MalformedLine(t *testing.T) {
	// Lines with no colon or too few fields are silently skipped.
	profile := "mode: atomic\nnot-a-valid-line\nexample.com/a.go:1.1,2.1 1 1\n"
	tmp := t.TempDir()
	p := filepath.Join(tmp, "c.out")
	if err := os.WriteFile(p, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}
	pct, err := computeCoverage(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if pct != 100.0 {
		t.Errorf("pct = %v, want 100", pct)
	}
}

func TestCoverage_ExcludePattern(t *testing.T) {
	// Profile: a.go has 2 stmts covered, b.go has 2 stmts NOT covered.
	// Excluding b.go should give 100% (only a.go counts).
	// Not excluding should give 50%.
	profile := "mode: atomic\nexample.com/pkg/a.go:1.1,2.1 2 1\nexample.com/pkg/b.go:3.1,4.1 2 0\n"
	tmp := t.TempDir()
	p := filepath.Join(tmp, "c.out")
	if err := os.WriteFile(p, []byte(profile), 0o600); err != nil {
		t.Fatal(err)
	}

	// Without exclude: 50%
	pct, err := computeCoverage(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if pct != 50.0 {
		t.Errorf("no-exclude pct = %v, want 50", pct)
	}

	// Excluding b.go: only a.go counts → 100%
	pct, err = computeCoverage(p, []string{"pkg/b.go"})
	if err != nil {
		t.Fatal(err)
	}
	if pct != 100.0 {
		t.Errorf("exclude pct = %v, want 100", pct)
	}
}
