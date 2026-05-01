// Package config owns the server's YAML schema, defaults, and the merge
// of file values plus CLI flag overrides. There are no environment
// variable reads in this package by design.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level server configuration.
type Config struct {
	HTTP    HTTPConfig `yaml:"http"`
	Auth    AuthConfig `yaml:"auth"`
	DataDir string     `yaml:"data_dir"`
}

// HTTPConfig holds HTTP server settings.
type HTTPConfig struct {
	Address        string    `yaml:"address"`
	TLS            TLSConfig `yaml:"tls"`
	TrustedProxies []string  `yaml:"trusted_proxies"`
	CORSOrigin     string    `yaml:"cors_origin"`
	HSTSEnabled    bool      `yaml:"hsts_enabled"`
}

// TLSConfig holds TLS settings.
type TLSConfig struct {
	Enabled   bool   `yaml:"enabled"`
	CertFile  string `yaml:"cert_file"`
	KeyFile   string `yaml:"key_file"`
	ChainFile string `yaml:"chain_file"`
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	MaxFailedAttemptsBeforeLockout int `yaml:"max_failed_attempts_before_lockout"`
	RateLimitWindowMinutes         int `yaml:"rate_limit_window_minutes"`
	RateLimitMaxAttempts           int `yaml:"rate_limit_max_attempts"`
	SessionLifetimeHours           int `yaml:"session_lifetime_hours"`
}

// CLIFlags holds CLI-supplied overrides plus a Set bit per field.
type CLIFlags struct {
	ConfigFile    string
	ConfigFileSet bool

	HTTPAddr    string
	HTTPAddrSet bool

	TLSEnabled    bool
	TLSEnabledSet bool
	TLSCertFile   string
	TLSCertSet    bool
	TLSKeyFile    string
	TLSKeySet     bool
	TLSChainFile  string
	TLSChainSet   bool

	DataDir    string
	DataDirSet bool
}

// defaultConfig returns the hard-coded defaults.
func defaultConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Address: ":<HTTP_PORT>",
			TLS:     TLSConfig{Enabled: false},
		},
		Auth: AuthConfig{
			MaxFailedAttemptsBeforeLockout: 10,
			RateLimitWindowMinutes:         15,
			RateLimitMaxAttempts:           10,
			SessionLifetimeHours:           24,
		},
	}
}

// LoadConfig builds a Config by:
//  1. starting with defaults,
//  2. merging the file at configPath if present (or erroring if --config was
//     explicit and the file is missing/invalid),
//  3. applying CLI overrides,
//  4. validating.
func LoadConfig(configPath string, flags CLIFlags) (*Config, error) {
	cfg := defaultConfig()

	if configPath != "" {
		fileCfg, err := loadConfigFile(configPath)
		switch {
		case err == nil:
			mergeConfig(cfg, fileCfg)
		case flags.ConfigFileSet:
			return nil, fmt.Errorf("loading config %s: %w", configPath, err)
		default:
			// Default-path file missing/invalid: silently use defaults.
		}
	}

	applyCLIFlags(cfg, flags)

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return cfg, nil
}

func loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path) // #nosec G304 - administrator-supplied path
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}
	return &cfg, nil
}

// mergeConfig copies non-zero fields from src into dest. Boolean fields are
// only copied when they are true so a config file that omits hsts_enabled
// does not flip an existing true to false.
func mergeConfig(dest, src *Config) {
	if src.HTTP.Address != "" {
		dest.HTTP.Address = src.HTTP.Address
	}
	if src.HTTP.TLS.Enabled {
		dest.HTTP.TLS.Enabled = true
	}
	if src.HTTP.TLS.CertFile != "" {
		dest.HTTP.TLS.CertFile = src.HTTP.TLS.CertFile
	}
	if src.HTTP.TLS.KeyFile != "" {
		dest.HTTP.TLS.KeyFile = src.HTTP.TLS.KeyFile
	}
	if src.HTTP.TLS.ChainFile != "" {
		dest.HTTP.TLS.ChainFile = src.HTTP.TLS.ChainFile
	}
	if len(src.HTTP.TrustedProxies) > 0 {
		dest.HTTP.TrustedProxies = src.HTTP.TrustedProxies
	}
	if src.HTTP.CORSOrigin != "" {
		dest.HTTP.CORSOrigin = src.HTTP.CORSOrigin
	}
	if src.HTTP.HSTSEnabled {
		dest.HTTP.HSTSEnabled = true
	}
	if src.Auth.MaxFailedAttemptsBeforeLockout > 0 {
		dest.Auth.MaxFailedAttemptsBeforeLockout = src.Auth.MaxFailedAttemptsBeforeLockout
	}
	if src.Auth.RateLimitWindowMinutes > 0 {
		dest.Auth.RateLimitWindowMinutes = src.Auth.RateLimitWindowMinutes
	}
	if src.Auth.RateLimitMaxAttempts > 0 {
		dest.Auth.RateLimitMaxAttempts = src.Auth.RateLimitMaxAttempts
	}
	if src.Auth.SessionLifetimeHours > 0 {
		dest.Auth.SessionLifetimeHours = src.Auth.SessionLifetimeHours
	}
	if src.DataDir != "" {
		dest.DataDir = src.DataDir
	}
}

func applyCLIFlags(cfg *Config, f CLIFlags) {
	if f.HTTPAddrSet {
		cfg.HTTP.Address = f.HTTPAddr
	}
	if f.TLSEnabledSet {
		cfg.HTTP.TLS.Enabled = f.TLSEnabled
	}
	if f.TLSCertSet {
		cfg.HTTP.TLS.CertFile = f.TLSCertFile
	}
	if f.TLSKeySet {
		cfg.HTTP.TLS.KeyFile = f.TLSKeyFile
	}
	if f.TLSChainSet {
		cfg.HTTP.TLS.ChainFile = f.TLSChainFile
	}
	if f.DataDirSet {
		cfg.DataDir = f.DataDir
	}
}

func validateConfig(cfg *Config) error {
	if cfg.HTTP.TLS.Enabled {
		if cfg.HTTP.TLS.CertFile == "" {
			return fmt.Errorf("tls cert file is required when tls.enabled is true")
		}
		if cfg.HTTP.TLS.KeyFile == "" {
			return fmt.Errorf("tls key file is required when tls.enabled is true")
		}
	}
	return nil
}
