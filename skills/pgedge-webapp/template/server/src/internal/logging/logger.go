// Package logging wraps log/slog so the rest of the server can take a
// *slog.Logger by injection. JSON by default, text+debug when debug=true.
package logging

import (
	"io"
	"log/slog"
)

// New returns a *slog.Logger writing to w. When debug is true the level is
// set to debug and a text handler is used (easier on humans during dev).
// Otherwise the level is info and a JSON handler is used.
func New(w io.Writer, debug bool) *slog.Logger {
	if debug {
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
