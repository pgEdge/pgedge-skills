#!/usr/bin/env bash
#
# Per-task test runner used during development of the pgedge-webapp skill.
# Copies template/server into a tempdir, substitutes placeholders, runs go
# test against the supplied package selector. Intended only as a stand-in
# until the full harness (test-scaffold.sh) can run end-to-end (which needs
# the genhash tool from Phase 18).
#
# Usage: bash skills/pgedge-webapp/tools/scratch-test.sh ./internal/auth/...
#
# Will be removed once the main harness is functional.
#
set -euo pipefail

# Force auto-toolchain so the script works on hosts that don't yet ship
# Go 1.26. Override only with explicit HARNESS_GOTOOLCHAIN if needed.
export GOTOOLCHAIN="${HARNESS_GOTOOLCHAIN:-auto}"

SELECTOR="${1:-./...}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
TEMPLATE_SERVER="$REPO_ROOT/skills/pgedge-webapp/template/server"

SCRATCH="$(mktemp -d -t pgedge-scratch-XXXXXX)"
trap 'rm -rf "$SCRATCH"' EXIT

cp -R "$TEMPLATE_SERVER/." "$SCRATCH/"

# Path-based placeholders
find "$SCRATCH" -depth -name '*<binary>*' | while read -r path; do
    mv "$path" "${path//<binary>/scratch-server}"
done

# Content placeholders
python3 - "$SCRATCH" <<'PYEOF'
import os, pathlib, sys
target = pathlib.Path(sys.argv[1])
subs = {
    "<MODULE_PATH>":     "example.com/scratch",
    "<BINARY_NAME>":     "scratch-server",
    "<BINARY_GO_IDENT>": "scratch_server",
    "<HTTP_PORT>":       "8080",
    "<COOKIE_NAME>":     "scratch_server_session",
    "<PROJECT_NAME>":    "Scratch App",
    "<PROJECT_SLUG>":    "scratch-app",
    "<CURRENT_YEAR>":    "2026",
    "<GITHUB_REPO>":     "example/scratch-app",
}
binary_extensions = {".png", ".ico", ".jpg", ".jpeg", ".gif"}
for path in target.rglob("*"):
    if not path.is_file() or path.suffix.lower() in binary_extensions:
        continue
    try:
        text = path.read_text()
    except UnicodeDecodeError:
        continue
    new = text
    for key, val in subs.items():
        new = new.replace(key, val)
    if new != text:
        path.write_text(new)
PYEOF

cd "$SCRATCH/src"
go mod tidy
exec go test -race "$SELECTOR"
