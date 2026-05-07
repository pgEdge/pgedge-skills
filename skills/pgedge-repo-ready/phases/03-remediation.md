# Phase 3 — Remediation

Walk through each failing or partial check and apply fixes. Use
templates from the `templates/` directory with placeholder
substitution per `PLACEHOLDERS.md`.

## Remediation Modes

### Batch (mode=batch)

Process all non-passing checks in tier order (1 through 6).
For each check:

1. State what will be created or changed (one sentence)
2. Apply the fix
3. Verify if applicable
4. Move to next check

Do not pause for confirmation between checks. Report a summary
at the end.

### Selected (mode=selected)

Same as batch, but only process checks in the user-selected
tiers.

### Interactive (mode=interactive)

Process non-passing checks in tier order. For each check:

1. Explain what will be created or changed
2. Show a preview of the content (first 20 lines of template, or
   full content if short)
3. Ask the user to confirm: [Y]es / [S]kip / [E]dit first
4. If Yes: apply and verify
5. If Skip: move to next check
6. If Edit: let the user describe changes, modify the content,
   then apply

## Placeholder Resolution

Before applying any templates, resolve all placeholders. Run
detection in this order:

```bash
# PROJECT_NAME
if [ -f go.mod ]; then
  grep "^module" go.mod | awk '{print $2}' | rev \
    | cut -d/ -f1 | rev
elif [ -f pyproject.toml ]; then
  grep "^name" pyproject.toml | head -1 \
    | sed 's/.*= *"\(.*\)"/\1/'
elif [ -f package.json ]; then
  grep '"name"' package.json | head -1 \
    | sed 's/.*: *"\(.*\)".*/\1/'
else
  basename $(pwd)
fi

# GITHUB_REPO
GITHUB_REPO=$(git remote get-url origin 2>/dev/null \
  | sed 's/.*github.com[:/]\(.*\)\.git/\1/' \
  | sed 's/.*github.com[:/]\(.*\)/\1/')
# If empty (no remote), ask the user for org/repo
if [ -z "$GITHUB_REPO" ]; then
  echo "No git remote found — ask user for org/repo"
fi

# GITHUB_ORG
echo $GITHUB_REPO | cut -d/ -f1

# CURRENT_YEAR
date +%Y

# COPYRIGHT_RANGE (check LICENSE.md for existing start year)
if [ -f LICENSE.md ]; then
  grep -oE "20[0-9]{2}" LICENSE.md | head -1
fi
# If found: "<start_year>-<CURRENT_YEAR>"
# If not: "<CURRENT_YEAR>"

# PRIMARY_OWNER — ask the user
# GO_VERSION, NODE_VERSION, PYTHON_VERSION — from project files
# PG_VERSIONS — always [16, 17, 18]
```

Confirm resolved values with the user before proceeding:

```
Resolved placeholders:
  PROJECT_NAME:    mm-ready-go
  GITHUB_REPO:     pgEdge/mm-ready-go
  GITHUB_ORG:      pgEdge
  CURRENT_YEAR:    2026
  COPYRIGHT_RANGE: 2025-2026
  PRIMARY_OWNER:   ?  (enter GitHub username)
  GO_VERSION:      1.26
  PG_VERSIONS:     [16, 17, 18]

Correct? [Y/n]
```

## Template Application

For each template-backed check:

1. Determine the source template path (relative to
   `templates/` in this skill)
2. Determine the target path in the repo
3. Read the template content
4. Substitute all applicable placeholders
5. Write the file to the target path
6. Verify the result (see verification section)

### Target Path Mapping

| Template | Target |
|----------|--------|
| `github/CODEOWNERS` | `.github/CODEOWNERS` |
| `github/SECURITY.md` | `.github/SECURITY.md` |
| `github/CONTRIBUTING.md` | `CONTRIBUTING.md` |
| `github/ISSUE_TEMPLATE/*` | `.github/ISSUE_TEMPLATE/*` |
| `github/PULL_REQUEST_TEMPLATE.md` | `.github/PULL_REQUEST_TEMPLATE.md` |
| `github/dependabot.yml` | `.github/dependabot.yml` |
| `ci/ci-go.yml` | `.github/workflows/ci.yml` (or `ci-<component>.yml` for multi-component) |
| `ci/ci-python.yml` | `.github/workflows/ci.yml` |
| `ci/ci-typescript.yml` | `.github/workflows/ci-client.yml` |
| `ci/ci-docs.yml` | `.github/workflows/ci-docs.yml` |
| `ci/ci-docker.yml` | `.github/workflows/ci-docker.yml` |
| `ci/release-go.yml` | `.github/workflows/release.yml` |
| `review/.coderabbit.yaml` | `.coderabbit.yaml` |
| `review/.codacy.yaml` | `.codacy.yaml` |
| `linting/.golangci.yml` | `.golangci.yml` (or per-component) |
| `linting/.pre-commit-config-base.yaml` | `.pre-commit-config.yaml` (when no language detected) |
| `linting/.pre-commit-config-{lang}.yaml` | `.pre-commit-config.yaml` |
| `linting/ruff.toml` | `ruff.toml` |
| `docs/mkdocs.yml` | `mkdocs.yml` |
| `docs/docs/*` | `docs/*` |
| `license/LICENSE.md` | `LICENSE.md` |
| `claude/CLAUDE.local.md` | `CLAUDE.local.md` |
| `claude/agents/*` | `.claude/agents/*` |

### Pre-Commit Config Selection

When **no languages** are detected, use
`.pre-commit-config-base.yaml` — base hooks (trailing-whitespace,
end-of-file-fixer, check-yaml, check-added-large-files) plus
gitleaks.

When **one language** is detected, use
`.pre-commit-config-{lang}.yaml` directly.

When **multiple languages** are detected, combine hooks from all
relevant `.pre-commit-config-{lang}.yaml` templates into a
single `.pre-commit-config.yaml`. Merge the `repos:` lists,
keeping one copy of shared repos (pre-commit-hooks, gitleaks)
and adding language-specific repos from each template.

### CI Workflow Selection

Apply CI workflow templates based on detected characteristics:

- `ci-go.yml` — if Go detected
- `ci-python.yml` — if Python detected
- `ci-typescript.yml` — if TypeScript detected
- `ci-docs.yml` — if `mkdocs.yml` exists or `docs/` directory
  present (regardless of language)
- `ci-docker.yml` — if `docker-compose.yml` or `Dockerfile`
  present
- `release-go.yml` — if Go binary project detected

### Multi-Component Repos

For repos with multiple components (e.g., server/, collector/,
client/), some templates need per-component application:

- `.golangci.yml` — copy to each Go component directory
- CI workflows — create one per component
  (ci-server.yml, ci-collector.yml, etc.)
- Makefile — check root-level orchestrator exists

### Existing File Handling

- If target file does not exist: create it from template
- If target file exists and check status is PARTIAL: show a diff
  of what would change, let the user decide
- If target file exists and check status is PASS: skip (do not
  overwrite passing files)
- Never silently overwrite existing files

## Detection-Driven Remediation

Not all checks are fixed by copying a template. The following
checks require generating content from detection results. These
are equally auto-fixable — do not skip them or mark them as
"needs human input."

### Q-07: Release Automation

If a Go binary entry point is detected (`func main()` in
`cmd/*/main.go`), apply `templates/ci/release-go.yml` to
`.github/workflows/release.yml`. Add `./cmd/<name>` as the
build target in the `go build` step. If no binary is detected,
SKIP this check.

### Q-08: License Scanning

If Go is detected, add a `licenses` target to the Makefile:

```makefile
licenses:
	@if ! command -v go-licenses &> /dev/null; then \
		go install github.com/google/go-licenses@latest; \
	fi
	go-licenses report ./... 2>/dev/null | tee NOTICE.txt
```

Add `licenses` to the `.PHONY` line. Add `NOTICE.txt` to
`.gitignore`.

For Python, add `pip-licenses` to dev dependencies and a
`licenses` Makefile target that runs `pip-licenses --format=csv
--output-file=NOTICE.txt`.

### D-03: Doc Styling

Find all markdown files in `docs/` with lines exceeding 79
characters (excluding code blocks, tables, URLs, and headings):

```bash
find docs/ -name "*.md" -exec \
  awk '/^```/,/^```/{next} /^\|/{next} length > 79 \
    {print FILENAME ":" NR}' {} \;
```

Reflow each long line to 79-character wrapping, preserving:

- Code blocks (fenced with ```)
- Tables (lines starting with `|`)
- URLs and links
- Headings

This is a mechanical transformation. Do not skip it.

### D-04: API Documentation Completeness

If the project exposes an API (has HTTP route handlers), compare
documented endpoints against actual routes:

1. Extract routes from the router file (look for chi, mux, gin,
   echo, or fiber route registrations)
2. Extract documented endpoints from API docs (grep for HTTP
   method + path patterns)
3. List undocumented routes grouped by resource

For each undocumented group, add a stub section to the API doc
with the method, path, and a `<!-- REVIEW: describe -->` marker.
The user fills in request/response details later, but the
skeleton must be complete. An API doc that lists 20 of 80
endpoints is worse than no API doc — it implies coverage it
does not have.

### D-05: Example Configs

If the project has configuration files (`.env.example`,
`config/*.example`, or documented environment variables), create
an `examples/` directory with copies of the example configs.
If no example configs exist but the project reads environment
variables, generate an `examples/.env.example` listing
detected variables with placeholder values.

### G-05: Dependabot Ecosystem Detection

The `templates/github/dependabot.yml` template only includes
`github-actions`. When applying it, add ecosystem entries
based on detected languages and infrastructure:

- Go detected: add `gomod` ecosystem, directory `/`
- Node/TypeScript detected: add `npm` ecosystem, directory
  set to the client path (e.g., `/client` or `/`)
- Python detected: add `pip` ecosystem, directory `/`
- Dockerfile present: add `docker` ecosystem, directory `/`

Each ecosystem entry uses `interval: "weekly"`. Example for
a Go + TypeScript repo:

```yaml
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"

  - package-ecosystem: "npm"
    directory: "/client"
    schedule:
      interval: "weekly"
```

### Q-05: Coverage Threshold Enforcement

When CI workflows exist (Q-04 is PASS) but lack coverage
threshold enforcement (Q-05 is PARTIAL or MISSING), patch
the existing workflow — do not replace it.

For Go, add a step after the test step that parses
`coverage.out` and fails if below 90%:

```yaml
- name: Check coverage threshold
  run: |
    COVERAGE=$(go tool cover -func=coverage.out \
      | grep total | awk '{print $3}' | tr -d '%')
    echo "Total coverage: ${COVERAGE}%"
    if (( $(echo "$COVERAGE < 90.0" | bc -l) )); then
      echo "Coverage ${COVERAGE}% is below 90% threshold"
      exit 1
    fi
```

For TypeScript (vitest), add `--coverage.thresholds.lines=90`
to the test command or configure in `vitest.config.ts`:

```typescript
coverage: {
  thresholds: { lines: 90, branches: 90 }
}
```

For Python, add `--fail-under=75` to the pytest-cov command.

This is a mechanical patch. Do not skip it.

### P-08: CI Postgres Version Matrix

When CI workflows exist but the Postgres version matrix is
incomplete, patch the existing workflow to include all three
required versions: 16, 17, and 18.

Find the `postgres-version` array in the matrix strategy
and ensure it contains `['16', '17', '18']`.

Also check the Postgres service image in CI workflows. If it
uses the upstream `postgres:` image instead of the pgEdge
image, update it to:

```yaml
image: ghcr.io/pgedge/pgedge-postgres:${{ matrix.postgres-version }}-spock5-standard
```

Using the upstream image in CI defeats the purpose of testing
against pgEdge's distribution. The CI image should match the
Docker Compose image (P-06).

### G-09: .gitignore

If `.gitignore` is MISSING, create it. If PARTIAL (exists
but missing key patterns), append the missing entries. Never
remove existing entries.

Start with the base patterns that apply to every repo:

```gitignore
# IDE and editor files
.idea/
.vscode/
*.swp
*.swo
*~
.DS_Store

# Secrets
.env
.env.local
.env.*.local
*.password

# Claude
.claude/settings.local.json
CLAUDE.local.md

# Audit artifacts
REPO_READY_AUDIT.md
NOTICE.txt
```

Then add language-specific patterns based on detection:

**Go detected:**

```gitignore
# Go
bin/
*.exe
*.dll
*.so
*.dylib
*.test
*.out
coverage.out
go.work
```

**Python detected:**

```gitignore
# Python
__pycache__/
*.py[cod]
*.egg-info/
dist/
build/
.venv/
venv/
.mypy_cache/
.ruff_cache/
htmlcov/
.coverage
```

**TypeScript/Node detected:**

```gitignore
# Node
node_modules/
dist/
*.tsbuildinfo
coverage/
```

**Docker detected:**

```gitignore
# Docker data
data/
```

When appending to an existing file, add a blank line before
the new section. Do not duplicate patterns already present.

### X-09: .devcontainer

Generate `.devcontainer/devcontainer.json` from detected
languages and services:

```json
{
  "name": "<PROJECT_NAME>",
  "features": {}
}
```

Populate `features` based on detection:

- Go detected: add `ghcr.io/devcontainers/features/go:1`
  with the detected Go version
- Node/TypeScript detected: add
  `ghcr.io/devcontainers/features/node:1` with detected
  Node version
- Python detected: add
  `ghcr.io/devcontainers/features/python:1` with detected
  Python version

If `docker-compose.yml` exists, add `dockerComposeFile` and
`service` fields. Add `forwardPorts` for detected services
(5432 for Postgres, 8080/8081 for API, 5173/3000 for frontend).

Add VS Code extensions based on detected languages:

- Go: `golang.go`
- TypeScript: `dbaeumer.vscode-eslint`, `esbenp.prettier-vscode`
- Python: `ms-python.python`, `charliermarsh.ruff`

## Verification

After applying templates, verify where possible:

```bash
# YAML validity (for all .yml/.yaml files created)
# NOTE: Skip mkdocs.yml — its !!python/name tags require
# pymdownx to be installed and will always fail safe_load.
# Validate mkdocs.yml with mkdocs build instead.
python3 -c "import yaml; yaml.safe_load(open('<file>'))" \
  2>/dev/null && echo "YAML valid" || echo "YAML invalid"

# MkDocs build (if docs templates were applied)
pip install mkdocs-material 2>/dev/null \
  && mkdocs build --strict 2>&1 | tail -5

# golangci-lint config (if Go linting was applied)
golangci-lint config verify 2>/dev/null

# Pre-commit hooks (if pre-commit was applied)
pre-commit validate-config 2>/dev/null
```

Report any verification failures and offer to fix.

## After Remediation

Check if CLAUDE.md needs attention (X-01 status is MISSING or
PARTIAL). If so, proceed to Phase 4. Otherwise, print a summary:

```
Remediation complete.
  Applied: NN fixes
  Skipped: NN checks
  Failed:  NN verifications

Updated REPO_READY_AUDIT.md with post-remediation status.
```

If Phase 4 is needed:

```
CLAUDE.md needs setup. Starting interview...
```

Proceed to Phase 4.
