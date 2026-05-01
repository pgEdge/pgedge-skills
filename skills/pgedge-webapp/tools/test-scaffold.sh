#!/usr/bin/env bash
#
# Scaffolds a project from skills/pgedge-webapp/template into a temp dir,
# replaces placeholders, generates the seed admin bcrypt hash, then runs
# the scaffolded project's full test suite. Exits non-zero on any failure.
#
# Usage: bash skills/pgedge-webapp/tools/test-scaffold.sh
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
TEMPLATE_DIR="$REPO_ROOT/skills/pgedge-webapp/template"
TOOLS_DIR="$REPO_ROOT/skills/pgedge-webapp/tools"

# Test inputs
PROJECT_NAME="Harness App"
PROJECT_SLUG="harness-app"
BINARY_NAME="harness-server"
BINARY_GO_IDENT="harness_server"
MODULE_PATH="github.com/example/harness"
HTTP_PORT="8080"
CURRENT_YEAR="$(date +%Y)"
COOKIE_NAME="${BINARY_GO_IDENT}_session"
ADMIN_PASSWORD="TestAdminPassword123!"

# Tempdir under the OS tempdir so artifacts auto-clean
TARGET="$(mktemp -d -t pgedge-webapp-scaffold-XXXXXX)"
echo "Scaffolding into: $TARGET"

# Copy the template tree
cp -R "$TEMPLATE_DIR/." "$TARGET/"

# Rename path placeholders
find "$TARGET" -depth -name '*<binary>*' | while read -r path; do
    new="${path//<binary>/$BINARY_NAME}"
    mv "$path" "$new"
done

# Replace content placeholders. Use a portable sed invocation by writing a
# tiny Python-based replacer to avoid GNU/BSD sed differences.
PROJECT_NAME="$PROJECT_NAME" \
PROJECT_SLUG="$PROJECT_SLUG" \
BINARY_NAME="$BINARY_NAME" \
BINARY_GO_IDENT="$BINARY_GO_IDENT" \
MODULE_PATH="$MODULE_PATH" \
HTTP_PORT="$HTTP_PORT" \
CURRENT_YEAR="$CURRENT_YEAR" \
COOKIE_NAME="$COOKIE_NAME" \
python3 - "$TARGET" <<'PYEOF'
import os, sys, pathlib
target = pathlib.Path(sys.argv[1])
subs = {
    "<PROJECT_NAME>":     os.environ["PROJECT_NAME"],
    "<PROJECT_SLUG>":     os.environ["PROJECT_SLUG"],
    "<BINARY_NAME>":      os.environ["BINARY_NAME"],
    "<BINARY_GO_IDENT>":  os.environ["BINARY_GO_IDENT"],
    "<MODULE_PATH>":      os.environ["MODULE_PATH"],
    "<HTTP_PORT>":        os.environ["HTTP_PORT"],
    "<CURRENT_YEAR>":     os.environ["CURRENT_YEAR"],
    "<COOKIE_NAME>":      os.environ["COOKIE_NAME"],
}
binary_extensions = {".png", ".ico", ".jpg", ".jpeg", ".gif"}
for path in target.rglob("*"):
    if not path.is_file():
        continue
    if path.suffix.lower() in binary_extensions:
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

# Generate the bcrypt hash for the seed admin
HASH="$(cd "$TOOLS_DIR" && go run ./genhash "$ADMIN_PASSWORD")"
SEED_FILE="$TARGET/server/src/cmd/$BINARY_NAME/seed.go"
python3 - "$SEED_FILE" "$HASH" <<'PYEOF'
import sys, pathlib
seed = pathlib.Path(sys.argv[1])
hash_value = sys.argv[2]
text = seed.read_text()
text = text.replace("<SEEDED_ADMIN_PASSWORD_HASH>", hash_value)
seed.write_text(text)
PYEOF

# Build + test the scaffolded server
(
    cd "$TARGET/server/src"
    go mod tidy
    go build ./...
    go vet ./...
    go test -race ./...
    cd "$REPO_ROOT"
)

# Build + test the scaffolded server's coverage gate
(
    cd "$TARGET/server"
    make coverage-check
)

# Build + test the scaffolded client
(
    cd "$TARGET/client"
    npm ci
    npm run typecheck
    npm run lint
    npm run test:coverage
)

# Helm lint (skipped if helm not installed; non-fatal)
if command -v helm >/dev/null 2>&1; then
    helm lint "$TARGET/deploy/helm/$BINARY_NAME"
else
    echo "helm not installed; skipping helm lint"
fi

echo "Scaffolding harness PASS: $TARGET"
