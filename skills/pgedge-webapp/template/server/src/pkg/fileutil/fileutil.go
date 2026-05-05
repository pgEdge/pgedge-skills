// Package fileutil resolves config file paths and provides small file
// helpers used by the server.
package fileutil

import (
	"os"
	"path/filepath"
)

// ConfigSource labels the origin of a resolved config path.
const (
	SourceCLI      = "cli"
	SourceUser     = "user"
	SourceSystem   = "system"
	SourceDefaults = "defaults"
)

// GetDefaultConfigPath walks the search order and returns the first existing
// path along with a source label. If no file exists, it returns the system
// path and source "system" (the caller can decide whether to error).
//
// Search order:
//  1. $XDG_CONFIG_HOME/pgedge/<configFilename>
//  2. os.UserConfigDir()/pgedge/<configFilename>
//     (~/.config on Linux, ~/Library/Application Support on macOS,
//     %AppData% on Windows)
//  3. /etc/pgedge/<configFilename>
//
// CLI override (--config) is handled by the caller before this function.
func GetDefaultConfigPath(configFilename string) (path, source string) {
	// Honor XDG_CONFIG_HOME explicitly so the same env var works on macOS
	// and Windows (os.UserConfigDir does not consult it on those platforms).
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		candidate := filepath.Join(xdg, "pgedge", configFilename)
		if FileExists(candidate) {
			return candidate, SourceUser
		}
	}
	if userDir, err := os.UserConfigDir(); err == nil {
		candidate := filepath.Join(userDir, "pgedge", configFilename)
		if FileExists(candidate) {
			return candidate, SourceUser
		}
	}
	return filepath.Join("/etc/pgedge", configFilename), SourceSystem
}

// FileExists returns true if the path refers to an existing regular file.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
