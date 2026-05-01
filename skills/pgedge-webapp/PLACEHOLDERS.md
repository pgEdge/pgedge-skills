# Placeholder Inventory

The scaffolding workflow substitutes the following placeholders in
template file *contents*. Path-based placeholders are listed
separately.

## Content placeholders

- `<PROJECT_NAME>` — human-readable display name.
- `<PROJECT_SLUG>` — kebab-case slug used in package identifiers.
- `<BINARY_NAME>` — kebab-case binary name used in commands.
- `<BINARY_GO_IDENT>` — snake_case form of the binary name for Go
  identifiers and the session cookie prefix.
- `<MODULE_PATH>` — full Go module import path.
- `<HTTP_PORT>` — numeric HTTP listen port.
- `<CURRENT_YEAR>` — four-digit year for copyright lines.
- `<COOKIE_NAME>` — session cookie name; derived from
  `<BINARY_GO_IDENT>` + `_session`.
- `<SEEDED_ADMIN_PASSWORD_HASH>` — bcrypt hash of the admin password,
  generated at scaffold time by `tools/genhash`.

## Path placeholders

- Directory `server/src/cmd/<binary>/` — rename to the actual binary.
- Filename `examples/<binary>.yaml` — rename to the actual binary.
- Directory `deploy/helm/<binary>/` — rename to the actual binary.

The harness `tools/test-scaffold.sh` is the canonical reference for the
substitution logic. Any change to placeholder names must be applied to
the harness in the same commit.
