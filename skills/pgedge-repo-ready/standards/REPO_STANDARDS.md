# pgEdge Repository Standards

The canonical checklist for pgEdge open-source repositories. Every
check has an ID (used in audit reports), a tier, pass criteria,
and the template or action used to remediate.

## How Tiers Work

- **Tiers 1-4** apply to every repository
- **Tier 5** applies per detected language (Go, Python, TypeScript)
- **Tier 6** applies when Postgres interaction is detected
  (pgx/psycopg in deps, .sql files, or docker-compose with
  Postgres services)

Tiers also define remediation priority — governance first,
language-specific last.

---

## Tier 1 — Governance

Every repo, regardless of language or purpose.

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| G-01 | License | `LICENSE.md` at repo root with PostgreSQL License, pgEdge Inc copyright, current year range | `templates/license/LICENSE.md` |
| G-02 | CODEOWNERS | `.github/CODEOWNERS` with at least one owner for `*` path | `templates/github/CODEOWNERS` |
| G-03 | Issue templates | `.github/ISSUE_TEMPLATE/bug_report.md` and `feature_request.md` both present | `templates/github/ISSUE_TEMPLATE/` |
| G-04 | PR template | `.github/PULL_REQUEST_TEMPLATE.md` with checklist | `templates/github/PULL_REQUEST_TEMPLATE.md` |
| G-05 | Dependabot | `.github/dependabot.yml` configured for relevant ecosystems | `templates/github/dependabot.yml` |
| G-06 | Security policy | `.github/SECURITY.md` with vulnerability reporting process | `templates/github/SECURITY.md` |
| G-07 | Contributing guide | `CONTRIBUTING.md` at repo root | `templates/github/CONTRIBUTING.md` |
| G-08 | README | `README.md` exists with project description, CI badge, and quickstart section | Manual (audit flags, user writes) |
| G-09 | .gitignore | `.gitignore` covers IDE files (.idea/, .vscode/), build artifacts (bin/, dist/, build/), secrets (.env, *.password), and `.claude/` | Manual (audit flags gaps) |

---

## Tier 2 — Code Quality

Language-agnostic setup with language-specific tool configurations.

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| Q-01 | Pre-commit hooks | `.pre-commit-config.yaml` with gitleaks + language-appropriate linters (ruff for Python, golangci-lint for Go) | `templates/linting/.pre-commit-config-{lang}.yaml` |
| Q-02 | CodeRabbit | `.coderabbit.yaml` with auto-review enabled, language-appropriate tools, noise reduction | `templates/review/.coderabbit.yaml` |
| Q-03 | Codacy | `.codacy.yaml` with exclusions for .claude/, test fixtures, examples | `templates/review/.codacy.yaml` |
| Q-04 | CI workflows | `.github/workflows/` with lint + test + build per detected language | `templates/ci/ci-{lang}.yml` |
| Q-05 | Coverage threshold | >=90% enforced in CI for Go and TypeScript; >=75% for Python | Checked in CI workflow templates |
| Q-06 | Secret scanning | gitleaks enabled in both pre-commit hooks and CodeRabbit config | Checked in Q-01 and Q-02 |
| Q-07 | Release automation | `.github/workflows/release*.yml` if project produces binaries or deployable artifacts | `templates/ci/release-go.yml` |
| Q-08 | License scanning | go-licenses or equivalent configured; NOTICE.txt generated for releases | Manual (audit flags) |

---

## Tier 3 — Documentation

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| D-01 | MkDocs config | `mkdocs.yml` with Material theme, pgEdge branding (logos, copyright), light/dark toggle | `templates/docs/mkdocs.yml` |
| D-02 | Docs structure | `docs/index.md` and `docs/changelog.md` present at minimum | `templates/docs/docs/` |
| D-03 | Doc styling | Markdown files wrap at 79 characters; pgEdge product names capitalized correctly (pgEdge Cloud, Spock, LOLOR, ACE) | Manual (audit samples files) |
| D-04 | API documentation | OpenAPI spec and/or ReDoc integration if project exposes a REST/HTTP API | Manual (audit flags) |
| D-05 | Example configs | `examples/` directory if project is configurable (has config files, CLI flags, or environment variables) | Manual (audit flags) |

---

## Tier 4 — Developer Experience

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| X-01 | CLAUDE.md | Present at repo root, structured with sections for context, architecture, commands, testing, code style; references `.claude/agents/` if agents exist | Generated via interview (Phase 4) |
| X-02 | Claude agents | `.claude/agents/` with language/domain-appropriate sub-agent definitions | `templates/claude/agents/` |
| X-03 | CLAUDE.local.md | Template present at repo root (must be in .gitignore) | `templates/claude/CLAUDE.local.md` |
| X-04 | Claude specs dir | `.claude/specs/` directory exists | Created as empty dir |
| X-05 | Claude plans dir | `.claude/plans/` directory exists | Created as empty dir |
| X-06 | Makefile | Present with at minimum: `build`, `test`, `lint`, `clean` targets | Manual (audit flags missing targets) |
| X-07 | Docker Compose | `docker-compose.yml` present if project has services (skip for pure libraries/CLIs) | Manual (audit flags) |
| X-08 | Conventional commits | Last 20 commits follow conventional commit style (feat:, fix:, docs:, chore:, etc.) | Manual (audit samples git log) |
| X-09 | .devcontainer | `.devcontainer/devcontainer.json` configured for project languages and tools | Manual (audit flags) |

---

## Tier 5 — Language-Specific

### Go (detected if .go files or go.mod present)

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| L-GO-01 | golangci-lint config | `.golangci.yml` with v2 format, standard linters: errcheck, gosec, govet, ineffassign, misspell, staticcheck, unused | `templates/linting/.golangci.yml` |
| L-GO-02 | gofmt compliance | All .go files pass `gofmt -l .` (no output) | Checked via command |
| L-GO-03 | Race detector | CI test commands include `-race` flag | Checked in CI workflow |
| L-GO-04 | Cross-platform build | Makefile has targets for linux/darwin (amd64/arm64) if project produces a binary | Manual (audit flags) |

### Python (detected if .py files or pyproject.toml present)

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| L-PY-01 | ruff config | `ruff.toml` or `[tool.ruff]` in pyproject.toml with standard rules (E, W, F, I, UP, B, SIM, RUF, D), Google docstring convention | `templates/linting/ruff.toml` |
| L-PY-02 | Type checking | pyright configured in strict mode (pyproject.toml or pyrightconfig.json) | Checked in pyproject.toml |
| L-PY-03 | Docstring coverage | interrogate configured with >=80% threshold | Checked in pyproject.toml |
| L-PY-04 | pyproject.toml | Present with pgEdge author, PostgreSQL License reference, repository URL | Manual (audit checks fields) |

### TypeScript/React (detected if .ts/.tsx files or package.json with typescript dep)

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| L-TS-01 | ESLint config | `eslint.config.js` or equivalent present | Manual (audit flags) |
| L-TS-02 | Prettier | `.prettierrc` or prettier config in package.json | Manual (audit flags) |
| L-TS-03 | TypeScript strict | `strict: true` in `tsconfig.json` | Checked via read |
| L-TS-04 | Vitest | Vitest configured with @vitest/coverage-v8 | Checked in package.json |

---

## Tier 6 — Postgres

Detected if pgx/psycopg in dependencies, .sql files present, or
docker-compose references a Postgres service.

| ID | Check | Pass Criteria | Template |
|----|-------|---------------|----------|
| P-01 | pgx/v5 driver | `go.mod` imports `github.com/jackc/pgx/v5` (Go) or psycopg/psycopg2 in requirements (Python) | Manual (audit flags wrong driver) |
| P-02 | SQL naming | snake_case for all table names, column names, index names, constraint names in .sql files and embedded SQL | Manual (audit samples SQL) |
| P-03 | TIMESTAMPTZ | No bare `TIMESTAMP` (without time zone) in CREATE TABLE statements or schema definitions | Manual (audit greps for TIMESTAMP) |
| P-04 | Parameterized queries | No string concatenation or f-string/fmt.Sprintf in SQL query construction | Manual (audit greps for patterns) |
| P-05 | Pool configuration | Connection pool explicitly sets MaxConns, HealthCheckPeriod, MaxConnIdleTime (not using bare defaults) | Manual (audit greps pool config) |
| P-06 | Docker image | docker-compose AND CI workflows use `ghcr.io/pgedge/pgedge-postgres` image, not upstream `postgres` | Manual (audit checks compose and CI) |
| P-07 | Schema comments | `COMMENT ON` statements present for tables, columns, indexes, and constraints in schema definitions | Manual (audit samples SQL) |
| P-08 | CI PG version matrix | CI workflows test against PostgreSQL 16, 17, and 18 | Checked in CI workflow |

---

## Postgres Conventions Reference

For use during audit and CLAUDE.md generation. Derived from
control-plane, ai-dba-workbench, pgedge-postgres-mcp.

**Driver:** pgx/v5 (github.com/jackc/pgx/v5) for Go. psycopg2 or
psycopg for Python.

**Naming:**
- Tables: `snake_case` (e.g., `connections`, `cluster_groups`)
- Columns: `snake_case` (e.g., `created_at`, `owner_username`)
- Indexes: `idx_{table}_{column}` (e.g., `idx_connections_name`)
- Unique indexes: `idx_{table}_{column}_unique` or
  `{table}_{cols}_unique`
- Check constraints: `chk_{description}` (e.g., `chk_owner`)
- Foreign keys: `fk_{child}_{parent}` (e.g.,
  `fk_connections_cluster_id`)

**Timestamps:** Always `TIMESTAMPTZ`, never bare `TIMESTAMP`.
Default: `CURRENT_TIMESTAMP`.

**Primary keys:**
- `SERIAL` or `BIGSERIAL` for local identity columns
- `UUID` for distributed or workflow identifiers
- `GENERATED ALWAYS AS IDENTITY` is acceptable

**Parameterized queries:** Use `pgx.NamedArgs` with `@param`
syntax in Go. Use `%s` placeholders with psycopg in Python.
Never concatenate strings into SQL.

**Connection pooling:** `pgxpool.Pool` with explicit configuration:
- `MaxConns` — set based on expected load
- `HealthCheckPeriod` — 30s to 1m
- `MaxConnIdleTime` — set to prevent stale connections

**Error handling:** Use `github.com/jackc/pgerrcode` to classify
PostgreSQL errors by code.

**Docker image:** `ghcr.io/pgedge/pgedge-postgres:18-spock5-standard`

**Health checks:** `pg_isready -U postgres` in docker-compose
healthcheck.

**Schema documentation:** `COMMENT ON TABLE`, `COMMENT ON COLUMN`,
`COMMENT ON INDEX`, `COMMENT ON CONSTRAINT` for all objects.

**Migrations:** Embedded in application code (Go or Python).
Idempotent (`IF NOT EXISTS`). Version-tracked in a schema_version
table or equivalent.

**Config priority:** Explicit config struct > `PGEDGE_*`
environment variables > hardcoded defaults.

**CI matrix:** Test against PostgreSQL 16, 17, 18.
