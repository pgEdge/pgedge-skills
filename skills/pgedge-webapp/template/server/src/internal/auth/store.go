// Package auth owns the sqlite-backed user and session store, the session
// middleware, the rate limiter, and the HTTP handlers for /api/v1/auth.
package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"golang.org/x/crypto/bcrypt"
)

// User represents a row from the users table.
type User struct {
	ID               int64
	Username         string
	FullName         string
	Email            string
	IsSuperuser      bool
	Enabled          bool
	FailedLoginCount int
	LockedUntil      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Store wraps the sqlite handle. Construct with Open; release with Close.
type Store struct {
	db *sql.DB
}

const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    email TEXT,
    is_superuser INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1,
    failed_login_count INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    token_hash TEXT PRIMARY KEY,
    username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_sessions_username   ON sessions(username);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
`

// Open opens (or creates) the sqlite file at path and applies the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite",
		path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := db.Exec(schemaSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close releases the underlying database handle.
func (s *Store) Close() error { return s.db.Close() }

// CountUsers returns the number of rows in users.
func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

// ErrUserNotFound indicates a missing user row.
var ErrUserNotFound = errors.New("user not found")

// CreateUserParams collects the fields required to create a user.
type CreateUserParams struct {
	Username    string
	Password    string
	FullName    string
	Email       string
	IsSuperuser bool
}

// CreateUser inserts a user, hashing Password with bcrypt cost 12.
func (s *Store) CreateUser(ctx context.Context, p CreateUserParams) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (username, password_hash, full_name, email, is_superuser)
		VALUES (?, ?, ?, ?, ?)
	`, p.Username, string(hash), p.FullName, p.Email, p.IsSuperuser)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// CreateUserWithHash inserts a user when the bcrypt hash is supplied
// directly (used by the seed-on-first-run path).
func (s *Store) CreateUserWithHash(ctx context.Context, username, hash string, isSuperuser bool) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (username, password_hash, is_superuser) VALUES (?, ?, ?)`,
		username, hash, isSuperuser)
	return err
}

// GetUser returns the user row.
func (s *Store) GetUser(ctx context.Context, username string) (*User, error) {
	return scanUser(s.db.QueryRowContext(ctx, userSelect+` WHERE username = ?`, username))
}

// ListUsers returns all users ordered by username.
func (s *Store) ListUsers(ctx context.Context) ([]*User, error) {
	rows, err := s.db.QueryContext(ctx, userSelect+` ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// DeleteUser removes the user; sessions cascade-delete via the FK.
func (s *Store) DeleteUser(ctx context.Context, username string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE username = ?`, username)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

// SetEnabled toggles the enabled flag.
func (s *Store) SetEnabled(ctx context.Context, username string, enabled bool) error {
	return s.update(ctx, username, "enabled = ?", enabled)
}

// SetSuperuser toggles is_superuser.
func (s *Store) SetSuperuser(ctx context.Context, username string, super bool) error {
	return s.update(ctx, username, "is_superuser = ?", super)
}

// UpdateProfile updates full_name and email.
func (s *Store) UpdateProfile(ctx context.Context, username, fullName, email string) error {
	return s.update(ctx, username, "full_name = ?, email = ?", fullName, email)
}

// UpdatePassword re-hashes and stores a new password and clears lockout.
func (s *Store) UpdatePassword(ctx context.Context, username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	return s.update(ctx, username,
		"password_hash = ?, failed_login_count = 0, locked_until = NULL",
		string(hash))
}

func (s *Store) update(ctx context.Context, username, set string, args ...any) error {
	args = append(args, username)
	res, err := s.db.ExecContext(ctx,
		"UPDATE users SET "+set+", updated_at = CURRENT_TIMESTAMP WHERE username = ?",
		args...)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}

const userSelect = `
SELECT id, username, full_name, email, is_superuser, enabled,
       failed_login_count, locked_until, created_at, updated_at
FROM users`

type rowScanner interface{ Scan(...any) error }

func scanUser(r rowScanner) (*User, error) {
	var u User
	var locked sql.NullTime
	err := r.Scan(&u.ID, &u.Username, &u.FullName, &u.Email,
		&u.IsSuperuser, &u.Enabled, &u.FailedLoginCount, &locked,
		&u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if locked.Valid {
		u.LockedUntil = &locked.Time
	}
	return &u, nil
}

// LockoutConfig parameterizes VerifyPassword lockout behavior.
type LockoutConfig struct {
	MaxFailedAttempts int
	LockoutDuration   time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
	ErrAccountDisabled    = errors.New("account disabled")
)

// VerifyPassword checks the credentials and applies the lockout policy.
func (s *Store) VerifyPassword(ctx context.Context, username, password string, cfg LockoutConfig) (*User, error) {
	var (
		passwordHash string
		enabled      bool
		failedCount  int
		lockedUntil  sql.NullTime
	)
	err := s.db.QueryRowContext(ctx, `
		SELECT password_hash, enabled, failed_login_count, locked_until
		FROM users WHERE username = ?
	`, username).Scan(&passwordHash, &enabled, &failedCount, &lockedUntil)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, ErrAccountDisabled
	}
	if lockedUntil.Valid && time.Now().Before(lockedUntil.Time) {
		return nil, ErrAccountLocked
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, s.recordFailedLogin(ctx, username, failedCount, cfg)
	}
	if _, err := s.db.ExecContext(ctx, `
		UPDATE users SET failed_login_count = 0, locked_until = NULL,
		                 updated_at = CURRENT_TIMESTAMP
		WHERE username = ?
	`, username); err != nil {
		return nil, err
	}
	return s.GetUser(ctx, username)
}

func (s *Store) recordFailedLogin(ctx context.Context, username string, prev int, cfg LockoutConfig) error {
	next := prev + 1
	var lockedUntil any
	if cfg.MaxFailedAttempts > 0 && next >= cfg.MaxFailedAttempts {
		d := cfg.LockoutDuration
		if d <= 0 {
			d = 15 * time.Minute
		}
		lockedUntil = time.Now().Add(d)
	}
	_, err := s.db.ExecContext(ctx, `
		UPDATE users SET failed_login_count = ?, locked_until = ?,
		                 updated_at = CURRENT_TIMESTAMP
		WHERE username = ?
	`, next, lockedUntil, username)
	if err != nil {
		return err
	}
	return ErrInvalidCredentials
}
