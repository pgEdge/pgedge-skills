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
//  1. $XDG_CONFIG_HOME/pgedge/<configFilename> (or ~/.config/pgedge/...)
//  2. /etc/pgedge/<configFilename>
//
// CLI override (--config) is handled by the caller before this function.
func GetDefaultConfigPath(configFilename string) (path, source string) {
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
