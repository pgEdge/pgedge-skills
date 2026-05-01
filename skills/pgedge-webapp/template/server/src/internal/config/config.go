// Package config owns the server's YAML schema, defaults, and the merge
// of file values plus CLI flag overrides. There are no environment
// variable reads in this package by design.
package config

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
