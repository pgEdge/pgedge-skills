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

func TestValidate_TLSMissingKey(t *testing.T) {
	cfg := defaultConfig()
	cfg.HTTP.TLS.Enabled = true
	cfg.HTTP.TLS.CertFile = "/c"
	// KeyFile omitted
	if err := validateConfig(cfg); err == nil {
		t.Fatal("missing key should fail")
	}
}

func TestLoadConfig_EmptyPath_SilentDefaults(t *testing.T) {
	// configPath="" means skip file loading entirely.
	cfg, err := LoadConfig("", CLIFlags{})
	if err != nil {
		t.Fatalf("expected defaults: %v", err)
	}
	if cfg.HTTP.Address != ":<HTTP_PORT>" {
		t.Errorf("address = %q", cfg.HTTP.Address)
	}
}

func TestLoadConfig_ConfigFileSet_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "bad.yaml")
	if err := os.WriteFile(p, []byte(":\tbad: yaml:\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadConfig(p, CLIFlags{ConfigFile: p, ConfigFileSet: true})
	if err == nil {
		t.Fatal("expected error for invalid YAML with ConfigFileSet=true")
	}
}

func TestLoadConfig_TLSChainCLIFlag(t *testing.T) {
	cfg := defaultConfig()
	applyCLIFlags(cfg, CLIFlags{
		TLSChainSet:  true,
		TLSChainFile: "/chain.pem",
	})
	if cfg.HTTP.TLS.ChainFile != "/chain.pem" {
		t.Errorf("chain = %q", cfg.HTTP.TLS.ChainFile)
	}
}

func TestMergeConfig_TrustedProxies(t *testing.T) {
	dest := defaultConfig()
	src := &Config{
		HTTP: HTTPConfig{
			TrustedProxies: []string{"10.0.0.1"},
			CORSOrigin:     "https://example.com",
		},
	}
	mergeConfig(dest, src)
	if len(dest.HTTP.TrustedProxies) != 1 || dest.HTTP.TrustedProxies[0] != "10.0.0.1" {
		t.Errorf("trusted proxies = %v", dest.HTTP.TrustedProxies)
	}
	if dest.HTTP.CORSOrigin != "https://example.com" {
		t.Errorf("cors origin = %q", dest.HTTP.CORSOrigin)
	}
}

func TestMergeConfig_TLS(t *testing.T) {
	dest := defaultConfig()
	src := &Config{
		HTTP: HTTPConfig{
			TLS: TLSConfig{
				Enabled:   true,
				CertFile:  "/cert",
				KeyFile:   "/key",
				ChainFile: "/chain",
			},
		},
	}
	mergeConfig(dest, src)
	if !dest.HTTP.TLS.Enabled {
		t.Error("TLS enabled not merged")
	}
	if dest.HTTP.TLS.CertFile != "/cert" {
		t.Errorf("cert = %q", dest.HTTP.TLS.CertFile)
	}
}
