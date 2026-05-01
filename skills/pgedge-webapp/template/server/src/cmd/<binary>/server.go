package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"<MODULE_PATH>/server/src/internal/api"
	"<MODULE_PATH>/server/src/internal/auth"
	"<MODULE_PATH>/server/src/internal/config"
	"<MODULE_PATH>/server/src/internal/logging"
	"<MODULE_PATH>/server/src/internal/openapi"
	"<MODULE_PATH>/server/src/pkg/fileutil"
)

// runServer is the testable bootstrap. main.go is a thin wrapper around it.
func runServer(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	defaultPath, defaultSource := fileutil.GetDefaultConfigPath("<BINARY_NAME>.yaml")
	flags, err := ParseFlags(args, defaultPath)
	if err != nil {
		return err
	}

	if flags.GenerateOpenAPISpec {
		b, err := openapi.MarshalJSON()
		if err != nil {
			return fmt.Errorf("openapi marshal: %w", err)
		}
		if _, err := stdout.Write(append(b, '\n')); err != nil {
			return err
		}
		return nil
	}

	log := logging.New(stderr, flags.Debug)

	source := defaultSource
	if flags.IsSet("config") {
		source = fileutil.SourceCLI
	}
	if !fileutil.FileExists(flags.ConfigFile) {
		source = fileutil.SourceDefaults
	}
	log.Info("config_resolved", "path", flags.ConfigFile, "source", source)

	cfg, err := config.LoadConfig(flags.ConfigFile, flags.ToCLIFlags())
	if err != nil {
		return err
	}

	dataDir := cfg.DataDir
	if dataDir == "" {
		exe, _ := os.Executable()
		dataDir = filepath.Join(filepath.Dir(exe), "data")
	}
	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return fmt.Errorf("data dir: %w", err)
	}

	store, err := auth.Open(filepath.Join(dataDir, "auth.db"))
	if err != nil {
		return err
	}
	defer func() { _ = store.Close() }()

	if err := seedAdminIfEmpty(ctx, store, flags, log); err != nil {
		return err
	}

	if flags.HasUserCommand() {
		return runUserCommand(ctx, store, flags, stdout)
	}
	if flags.Init {
		log.Info("init_complete", "data_dir", dataDir)
		return nil
	}

	return serveHTTP(ctx, cfg, store, log)
}

func seedAdminIfEmpty(ctx context.Context, store *auth.Store, flags *Flags, log *slog.Logger) error {
	n, err := store.CountUsers(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	hash := SeededAdminPasswordHash
	if flags.AdminPasswordFile != "" {
		b, err := os.ReadFile(flags.AdminPasswordFile) // #nosec G304
		if err != nil {
			return fmt.Errorf("admin password file: %w", err)
		}
		return store.CreateUser(ctx, auth.CreateUserParams{
			Username:    SeededAdminUsername,
			Password:    strings.TrimRight(string(b), "\n\r"),
			IsSuperuser: true,
		})
	}
	if hash == "<SEEDED_ADMIN_PASSWORD_HASH>" {
		return errors.New("seed admin hash not configured; provide -admin-password-file")
	}
	if err := store.CreateUserWithHash(ctx, SeededAdminUsername, hash, true); err != nil {
		return err
	}
	log.Warn("seeded_admin", "username", SeededAdminUsername,
		"msg", "rotate this password after first login")
	return nil
}

func serveHTTP(ctx context.Context, cfg *config.Config, store *auth.Store, log *slog.Logger) error {
	cookieName := "<COOKIE_NAME>"
	sessionLifetime := time.Duration(cfg.Auth.SessionLifetimeHours) * time.Hour
	handlers := auth.NewHandlers(auth.HandlersConfig{
		Store: store, CookieName: cookieName,
		SessionLifetime: sessionLifetime,
		Lockout: auth.LockoutConfig{
			MaxFailedAttempts: cfg.Auth.MaxFailedAttemptsBeforeLockout,
			LockoutDuration:   15 * time.Minute,
		},
		SecureCookie: cfg.HTTP.TLS.Enabled,
	})
	mw := auth.NewMiddleware(store, cookieName)
	rl := auth.NewRateLimiter(cfg.Auth.RateLimitMaxAttempts,
		time.Duration(cfg.Auth.RateLimitWindowMinutes)*time.Minute)

	handler := api.NewRouter(api.Deps{
		Store: store, Handlers: handlers, Middleware: mw, LoginLimiter: rl,
	})

	srv := &http.Server{
		Addr:              cfg.HTTP.Address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serveCtx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go cleanupExpiredSessions(serveCtx, store, log)

	errCh := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.HTTP.Address, "tls", cfg.HTTP.TLS.Enabled)
		if cfg.HTTP.TLS.Enabled {
			errCh <- srv.ListenAndServeTLS(cfg.HTTP.TLS.CertFile, cfg.HTTP.TLS.KeyFile)
		} else {
			errCh <- srv.ListenAndServe()
		}
	}()

	select {
	case <-serveCtx.Done():
		log.Info("shutting_down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func cleanupExpiredSessions(ctx context.Context, store *auth.Store, log *slog.Logger) {
	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if n, err := store.DeleteExpiredSessions(ctx); err != nil {
				log.Warn("cleanup_expired", "err", err)
			} else if n > 0 {
				log.Info("cleanup_expired", "removed", n)
			}
		}
	}
}
