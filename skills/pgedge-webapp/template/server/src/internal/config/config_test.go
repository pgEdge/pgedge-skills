package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	c := defaultConfig()

	if c.HTTP.Address != ":<HTTP_PORT>" {
		t.Errorf("default address = %q", c.HTTP.Address)
	}
	if c.Auth.MaxFailedAttemptsBeforeLockout != 10 {
		t.Errorf("default lockout = %d", c.Auth.MaxFailedAttemptsBeforeLockout)
	}
	if c.Auth.RateLimitWindowMinutes != 15 {
		t.Errorf("default window = %d", c.Auth.RateLimitWindowMinutes)
	}
	if c.Auth.RateLimitMaxAttempts != 10 {
		t.Errorf("default attempts = %d", c.Auth.RateLimitMaxAttempts)
	}
	if c.Auth.SessionLifetimeHours != 24 {
		t.Errorf("default session lifetime = %d", c.Auth.SessionLifetimeHours)
	}
}

func TestLoadConfig_FileMerge(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "cfg.yaml")
	body := []byte(`
http:
  address: ":9999"
  hsts_enabled: true
auth:
  rate_limit_window_minutes: 30
data_dir: /var/lib/test
`)
	if err := os.WriteFile(p, body, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(p, CLIFlags{ConfigFile: p, ConfigFileSet: true})
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.HTTP.Address != ":9999" {
		t.Errorf("address = %q", cfg.HTTP.Address)
	}
	if !cfg.HTTP.HSTSEnabled {
		t.Error("hsts not merged")
	}
	if cfg.Auth.RateLimitWindowMinutes != 30 {
		t.Errorf("window = %d", cfg.Auth.RateLimitWindowMinutes)
	}
	if cfg.Auth.MaxFailedAttemptsBeforeLockout != 10 {
		t.Error("default lockout overwritten by zero merge")
	}
	if cfg.DataDir != "/var/lib/test" {
		t.Errorf("data dir = %q", cfg.DataDir)
	}
}

func TestLoadConfig_MissingExplicitFile(t *testing.T) {
	_, err := LoadConfig("/no/such/file.yaml", CLIFlags{
		ConfigFile: "/no/such/file.yaml", ConfigFileSet: true,
	})
	if err == nil {
		t.Fatal("expected error when explicit --config is missing")
	}
}

func TestLoadConfig_SilentDefaultsWhenNoFile(t *testing.T) {
	cfg, err := LoadConfig("/no/such/file.yaml", CLIFlags{})
	if err != nil {
		t.Fatalf("expected silent defaults, got %v", err)
	}
	if cfg.HTTP.Address != ":<HTTP_PORT>" {
		t.Errorf("default lost")
	}
}

func TestApplyCLIFlags(t *testing.T) {
	cfg := defaultConfig()
	applyCLIFlags(cfg, CLIFlags{
		HTTPAddrSet:   true,
		HTTPAddr:      ":7777",
		TLSEnabledSet: true,
		TLSEnabled:    true,
		TLSCertSet:    true,
		TLSCertFile:   "/tmp/c.pem",
		TLSKeySet:     true,
		TLSKeyFile:    "/tmp/k.pem",
		DataDirSet:    true,
		DataDir:       "/tmp/data",
	})
	if cfg.HTTP.Address != ":7777" {
		t.Errorf("address = %q", cfg.HTTP.Address)
	}
	if !cfg.HTTP.TLS.Enabled || cfg.HTTP.TLS.CertFile != "/tmp/c.pem" {
		t.Errorf("tls = %+v", cfg.HTTP.TLS)
	}
	if cfg.DataDir != "/tmp/data" {
		t.Errorf("data dir = %q", cfg.DataDir)
	}
}

func TestValidate_TLSRequiresCertAndKey(t *testing.T) {
	cfg := defaultConfig()
	cfg.HTTP.TLS.Enabled = true
	if err := validateConfig(cfg); err == nil {
		t.Fatal("missing cert/key should fail")
	}
	cfg.HTTP.TLS.CertFile = "/c"
	cfg.HTTP.TLS.KeyFile = "/k"
	if err := validateConfig(cfg); err != nil {
		t.Fatalf("valid TLS rejected: %v", err)
	}
}
