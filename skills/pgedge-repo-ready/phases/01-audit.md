# Phase 1 — Audit

Read `standards/REPO_STANDARDS.md` before starting. Every check
in that document must be evaluated.

## Step 1: Detect Repo Characteristics

Run these commands and record the results. This determines which
tiers and checks apply.

### Language Detection

```bash
# Go
ls go.mod 2>/dev/null && echo "GO_DETECTED=true"
find . -name "*.go" -not -path "./.claude/*" | head -1

# Python
ls pyproject.toml 2>/dev/null && echo "PYTHON_DETECTED=true"
find . -name "*.py" -not -path "./.claude/*" | head -1

# TypeScript
ls tsconfig.json 2>/dev/null && echo "TS_DETECTED=true"
find . -name "*.ts" -o -name "*.tsx" | head -1
```

### Project Type Detection

```bash
# Binary/CLI (has main package or entry point)
grep -r "func main()" --include="*.go" -l 2>/dev/null
grep -r "entry_points\|console_scripts" pyproject.toml 2>/dev/null
grep '"bin"' package.json 2>/dev/null

# Service (has docker-compose or Dockerfile)
ls docker-compose.yml docker-compose.yaml Dockerfile 2>/dev/null

# Library (has publishable package config but no main)
grep -l "module\|name" go.mod pyproject.toml package.json 2>/dev/null

# Multi-component (has subdirectories with their own go.mod or
# package.json)
find . -mindepth 2 -maxdepth 2 \
  \( -name "go.mod" -o -name "package.json" \) 2>/dev/null
```

### Postgres Detection

```bash
# Go dependency
grep "jackc/pgx" go.mod 2>/dev/null
grep "lib/pq" go.mod 2>/dev/null

# Python dependency
grep -i "psycopg\|asyncpg" \
  pyproject.toml requirements.txt 2>/dev/null

# SQL files
find . -name "*.sql" -not -path "./.claude/*" | head -5

# Docker Postgres service
grep -i "postgres" docker-compose.yml 2>/dev/null
```

Record results as:
- `LANGUAGES`: list of detected languages (Go, Python, TypeScript)
- `PROJECT_TYPE`: one of (binary, service, library, multi-component)
- `HAS_POSTGRES`: true/false
- `HAS_DOCKER`: true/false

## Step 2: Determine Applicable Checks

- Tiers 1-4: always apply
- Tier 5 Go: apply if Go detected
- Tier 5 Python: apply if Python detected
- Tier 5 TypeScript: apply if TypeScript detected
- Tier 6: apply if HAS_POSTGRES is true

Count total applicable checks. This is the denominator for
scoring.

## Step 3: Run Each Check

For each applicable check in REPO_STANDARDS.md, evaluate and
record:

- **Check ID** (e.g., G-01)
- **Status**: one of PASS, PARTIAL, MISSING, SKIP
- **Finding**: what was actually found (or "not present")
- **Expected**: what the standard requires
- **Remediation**: what action would fix it (template path or
  manual step)

### Status Definitions

- **PASS** — fully meets the standard
- **PARTIAL** — exists but does not fully conform (e.g., CLAUDE.md
  exists but has no agents section; .golangci.yml exists but
  uses v1 format; CI exists but is missing coverage threshold)
- **MISSING** — not present at all
- **SKIP** — not applicable (language not detected, project type
  does not require it)

### Governance Checks (Tier 1)

```bash
# G-01: License
test -f LICENSE.md && grep -q "PostgreSQL License" LICENSE.md \
  && grep -q "pgEdge" LICENSE.md && echo "PASS" || echo "MISSING"

# G-02: CODEOWNERS
test -f .github/CODEOWNERS && echo "PASS" || echo "MISSING"

# G-03: Issue templates
test -f .github/ISSUE_TEMPLATE/bug_report.md \
  && test -f .github/ISSUE_TEMPLATE/feature_request.md \
  && echo "PASS" || echo "MISSING"

# G-04: PR template
test -f .github/PULL_REQUEST_TEMPLATE.md \
  && echo "PASS" || echo "MISSING"

# G-05: Dependabot
test -f .github/dependabot.yml && echo "PASS" || echo "MISSING"

# G-06: Security policy
test -f .github/SECURITY.md && echo "PASS" || echo "MISSING"

# G-07: Contributing guide
test -f CONTRIBUTING.md && echo "PASS" || echo "MISSING"

# G-08: README
test -f README.md && echo "EXISTS" || echo "MISSING"
# If exists, check for CI badge and quickstart section:
grep -q "badge\|CI\|Build" README.md 2>/dev/null \
  && echo "HAS_BADGE" || echo "NO_BADGE"
grep -qi "quickstart\|quick start\|getting started\|installation\|usage" \
  README.md 2>/dev/null \
  && echo "HAS_QUICKSTART" || echo "NO_QUICKSTART"
# PASS if all three; PARTIAL if exists but missing badge or
# quickstart

# G-09: .gitignore
test -f .gitignore && echo "EXISTS" || echo "MISSING"
# If exists, check for key patterns:
grep -q "\.claude" .gitignore 2>/dev/null \
  && echo "HAS_CLAUDE" || echo "NO_CLAUDE"
grep -q "\.env" .gitignore 2>/dev/null \
  && echo "HAS_ENV" || echo "NO_ENV"
# PASS if covers IDE, build, secrets, .claude/; PARTIAL if
# exists but missing key patterns
```

### Code Quality Checks (Tier 2)

```bash
# Q-01: Pre-commit hooks
test -f .pre-commit-config.yaml && echo "EXISTS" || echo "MISSING"
grep -q "gitleaks" .pre-commit-config.yaml 2>/dev/null \
  && echo "HAS_GITLEAKS" || echo "NO_GITLEAKS"

# Q-02: CodeRabbit
test -f .coderabbit.yaml && echo "EXISTS" || echo "MISSING"

# Q-03: Codacy
test -f .codacy.yaml && echo "EXISTS" || echo "MISSING"

# Q-04: CI workflows
ls .github/workflows/ci*.yml 2>/dev/null

# Q-05: Coverage threshold — check CI files for coverage flags
grep -r "coverage\|--cov\|coverprofile" \
  .github/workflows/ 2>/dev/null

# Q-06: Secret scanning — covered by Q-01 and Q-02 gitleaks check

# Q-07: Release automation
ls .github/workflows/release*.yml 2>/dev/null

# Q-08: License scanning
grep -r "go-licenses\|license-checker\|NOTICE" \
  Makefile .github/workflows/ 2>/dev/null
```

### Documentation Checks (Tier 3)

```bash
# D-01: MkDocs config
test -f mkdocs.yml && echo "EXISTS" || echo "MISSING"
grep -q "material" mkdocs.yml 2>/dev/null \
  && echo "HAS_MATERIAL" || echo "NO_MATERIAL"

# D-02: Docs structure
test -f docs/index.md && echo "HAS_INDEX" || echo "NO_INDEX"
test -f docs/changelog.md \
  && echo "HAS_CHANGELOG" || echo "NO_CHANGELOG"

# D-03: Doc styling — sample first 5 markdown files for line
# length
find docs/ -name "*.md" 2>/dev/null | head -5 | while read f; do
  awk 'length > 79' "$f" | head -3
done

# D-04: API docs — check for OpenAPI spec
find . -name "openapi*" -o -name "swagger*" 2>/dev/null | head -3

# D-05: Example configs
test -d examples && echo "HAS_EXAMPLES" || echo "NO_EXAMPLES"
```

### Developer Experience Checks (Tier 4)

```bash
# X-01: CLAUDE.md
test -f CLAUDE.md && echo "EXISTS" || echo "MISSING"
# If exists, check for structure:
grep -c "^##" CLAUDE.md 2>/dev/null  # count sections

# X-02: Claude agents
test -d .claude/agents && echo "EXISTS" || echo "MISSING"
ls .claude/agents/*.md 2>/dev/null | wc -l

# X-03: CLAUDE.local.md
test -f CLAUDE.local.md && echo "EXISTS" || echo "MISSING"

# X-04-05: Claude specs/plans dirs
test -d .claude/specs && echo "SPECS_EXISTS" || echo "SPECS_MISSING"
test -d .claude/plans && echo "PLANS_EXISTS" || echo "PLANS_MISSING"

# X-06: Makefile
test -f Makefile && echo "EXISTS" || echo "MISSING"
grep -c "^[a-z].*:" Makefile 2>/dev/null  # count targets
# Check for required targets:
for target in build test lint clean; do
  grep -q "^${target}:" Makefile 2>/dev/null \
    && echo "${target}:FOUND" || echo "${target}:MISSING"
done

# X-07: Docker Compose (only if service project)
test -f docker-compose.yml -o -f docker-compose.yaml \
  && echo "EXISTS" || echo "MISSING"

# X-08: Conventional commits
git log --oneline -20 2>/dev/null | head -20
# Count lines matching conventional commit pattern:
git log --oneline -20 2>/dev/null \
  | grep -cE "^[a-f0-9]+ (feat|fix|docs|chore|test|refactor|ci|build|style|perf)(\(.*\))?:"

# X-09: .devcontainer
test -f .devcontainer/devcontainer.json \
  && echo "EXISTS" || echo "MISSING"
```

### Language-Specific Checks (Tier 5)

Run only for detected languages.

**Go:**
```bash
# L-GO-01: golangci-lint config
test -f .golangci.yml && echo "EXISTS" || echo "MISSING"
# Also check subdirectories for multi-component repos:
find . -maxdepth 2 -name ".golangci.yml" 2>/dev/null

# L-GO-02: gofmt compliance
gofmt -l . 2>/dev/null | grep -v vendor | head -10

# L-GO-03: Race detector in CI
grep -r "\-race" .github/workflows/ 2>/dev/null

# L-GO-04: Cross-platform build (if binary project)
grep -E "GOOS|GOARCH|darwin|linux|windows" Makefile 2>/dev/null
```

**Python:**
```bash
# L-PY-01: ruff config
test -f ruff.toml && echo "EXISTS" \
  || grep -q "\[tool.ruff\]" pyproject.toml 2>/dev/null \
  && echo "IN_PYPROJECT" || echo "MISSING"

# L-PY-02: pyright strict
grep -q "typeCheckingMode.*strict" pyproject.toml 2>/dev/null \
  && echo "STRICT" || echo "NOT_STRICT"

# L-PY-03: interrogate
grep -q "interrogate" pyproject.toml 2>/dev/null \
  && echo "CONFIGURED" || echo "MISSING"

# L-PY-04: pyproject.toml metadata
test -f pyproject.toml && echo "EXISTS" || echo "MISSING"
grep -q "pgEdge\|pgedge" pyproject.toml 2>/dev/null \
  && echo "HAS_PGEDGE_AUTHOR" || echo "NO_PGEDGE_AUTHOR"
```

**TypeScript:**
```bash
# L-TS-01: ESLint
find . -maxdepth 2 -name "eslint.config.*" \
  -o -name ".eslintrc*" 2>/dev/null | head -3

# L-TS-02: Prettier
test -f .prettierrc && echo "EXISTS" \
  || grep -q "prettier" package.json 2>/dev/null \
  && echo "IN_PACKAGE_JSON" || echo "MISSING"

# L-TS-03: TypeScript strict
grep -q '"strict": true' tsconfig.json 2>/dev/null \
  && echo "STRICT" || echo "NOT_STRICT"

# L-TS-04: Vitest
grep -q "vitest" package.json 2>/dev/null \
  && echo "HAS_VITEST" || echo "MISSING"
grep -q "coverage-v8\|coverage-istanbul" package.json 2>/dev/null \
  && echo "HAS_COVERAGE" || echo "NO_COVERAGE"
```

### Postgres Checks (Tier 6)

Run only if HAS_POSTGRES is true.

```bash
# P-01: pgx/v5 driver
grep "jackc/pgx/v5" go.mod 2>/dev/null \
  && echo "PASS" || echo "CHECK"
# Flag if using lib/pq:
grep "lib/pq" go.mod 2>/dev/null && echo "WARN: using lib/pq"

# P-02: SQL naming — sample .sql files for non-snake_case
find . -name "*.sql" -exec grep -l "[A-Z]" {} \; 2>/dev/null \
  | head -5

# P-03: TIMESTAMPTZ — find bare TIMESTAMP (not TIMESTAMPTZ)
find . -name "*.sql" -exec \
  grep -Pn "TIMESTAMP(?!TZ)" {} \; 2>/dev/null | head -10
grep -rPn "TIMESTAMP(?!TZ)" --include="*.go" 2>/dev/null \
  | head -10

# P-04: Parameterized queries — flag string interpolation in SQL
grep -rn 'fmt.Sprintf.*SELECT\|fmt.Sprintf.*INSERT\|fmt.Sprintf.*UPDATE\|fmt.Sprintf.*DELETE' \
  --include="*.go" 2>/dev/null | head -10
grep -rn 'f".*SELECT\|f".*INSERT\|f".*UPDATE\|f".*DELETE' \
  --include="*.py" 2>/dev/null | head -10

# P-05: Pool configuration
grep -rn "MaxConns\|HealthCheckPeriod\|MaxConnIdleTime" \
  --include="*.go" 2>/dev/null | head -10

# P-06: Docker image (compose AND CI)
grep -i "pgedge.*postgres\|postgres:" \
  docker-compose.yml 2>/dev/null
# Also check CI workflow service images:
grep -i "image:.*postgres" \
  .github/workflows/ci*.yml 2>/dev/null

# P-07: Schema comments
find . -name "*.sql" -exec grep -l "COMMENT ON" {} \; 2>/dev/null
grep -rn "COMMENT ON" --include="*.go" 2>/dev/null | head -5

# P-08: CI PG version matrix
grep -A5 "matrix" .github/workflows/ci*.yml 2>/dev/null \
  | grep -i "postgres\|pg"
```

## Step 4: Compile Results

Store all check results in memory for Phase 2. For each check,
record:

1. Check ID (G-01, Q-02, etc.)
2. Check name
3. Status (PASS / PARTIAL / MISSING / SKIP)
4. Finding (actual state observed)
5. Expected (from REPO_STANDARDS.md)
6. Remediation action (template path or manual description)

Calculate per-tier scores: `passing / applicable` for each tier.
Calculate overall score: `total passing / total applicable`.

Proceed to Phase 2.
