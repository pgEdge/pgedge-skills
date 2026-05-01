package config

import "testing"

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
