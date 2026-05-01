package main

import (
	"fmt"
	"os"
	"strings"
)

// ResolvePassword returns the value from the explicit flag if set,
// otherwise from the password file if non-empty, otherwise empty string.
func ResolvePassword(flagValue string, flagSet bool, file string) (string, error) {
	if flagSet {
		return flagValue, nil
	}
	if file == "" {
		return "", nil
	}
	data, err := os.ReadFile(file) // #nosec G304 - administrator-supplied path
	if err != nil {
		return "", fmt.Errorf("reading password file: %w", err)
	}
	return strings.TrimRight(string(data), "\n\r"), nil
}
