# pgedge-repo-ready

Audit a repository against pgEdge engineering standards and
fix what's missing. Covers governance, CI, documentation,
developer experience, language tooling, and Postgres
conventions — 51 checks across 6 tiers.

## What It Does

Run this skill against any repo to get a scored audit and
guided remediation. The skill detects languages, project
type, and Postgres usage, then evaluates every applicable
check.

```
pgEdge Repo Ready Audit — imagineer
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Detected: Go, TypeScript
Project type: service
Postgres: yes

Tier 1 — Governance           4/9  ██████░░░░░░░░  44%
Tier 2 — Code Quality         3/8  █████░░░░░░░░░  38%
Tier 3 — Documentation        0/5  ░░░░░░░░░░░░░░   0%
Tier 4 — Developer Experience 6/9  █████████░░░░░  67%
Tier 5 — Go                   2/4  ███████░░░░░░░  50%
Tier 5 — TypeScript           3/4  ██████████░░░░  75%
Tier 6 — Postgres             5/8  ████████░░░░░░  63%

Overall: 23/47 (49%)
```

After the audit, choose a remediation mode:

- **Apply all** — batch-fix everything in tier order
- **Pick categories** — fix selected tiers only
- **One at a time** — preview each fix before applying

## The 6 Tiers

### Tier 1 — Governance (G-01 through G-09)

LICENSE.md, CODEOWNERS, issue/PR templates, dependabot,
security policy, contributing guide, README, .gitignore.

### Tier 2 — Code Quality (Q-01 through Q-08)

Pre-commit hooks with gitleaks, CodeRabbit, Codacy, CI
workflows, coverage thresholds (90% Go/TS, 75% Python),
secret scanning, release automation, license scanning.

### Tier 3 — Documentation (D-01 through D-05)

MkDocs Material config with pgEdge branding, docs
structure, 79-character line wrapping, API documentation
completeness, example configs.

### Tier 4 — Developer Experience (X-01 through X-09)

CLAUDE.md, Claude agents, CLAUDE.local.md, .claude/specs
and .claude/plans directories, Makefile targets,
docker-compose, conventional commits, .devcontainer.

### Tier 5 — Language-Specific (12 checks)

**Go:** golangci-lint v2 config, gofmt, race detector in
CI, cross-platform build.

**Python:** ruff config, pyright strict mode, interrogate
docstring coverage, pyproject.toml metadata.

**TypeScript:** ESLint, Prettier, strict mode, Vitest with
coverage.

### Tier 6 — Postgres (P-01 through P-08)

pgx/v5 driver, snake_case naming, TIMESTAMPTZ (no bare
TIMESTAMP), parameterized queries, pool configuration,
pgEdge Docker image, schema comments, CI PG version matrix
(16, 17, 18).

## How It Works

The skill runs four phases:

1. **Audit** — detect repo characteristics, run all
   applicable checks, record results
2. **Report** — print terminal summary, write detailed
   `REPO_READY_AUDIT.md`
3. **Remediation** — apply templates and detection-driven
   fixes for every non-passing check
4. **CLAUDE.md Interview** — if CLAUDE.md is missing or
   incomplete, run a structured Q&A to generate it

### Templates

The `templates/` directory contains 35 files across 7
categories: CI workflows, Claude agents, documentation
scaffolding, GitHub community files, license, linting
configs, and code review configs. Templates use
`<PLACEHOLDER>` syntax — values are auto-detected from the
repo and confirmed before application.

### Detection-Driven Fixes

Not everything is a template copy. The skill also generates
content from what it finds in the repo:

- Dependabot ecosystems matched to detected languages
- Coverage threshold enforcement patched into existing CI
- Release workflows wired to detected binary entry points
- License scanning targets added to Makefiles
- Long doc lines reflowed to 79 characters
- API doc stubs generated from router analysis
- .devcontainer configs built from detected languages
- CI Postgres matrix expanded to include all required
  versions

## Triggering the Skill

The skill activates when you mention repo readiness,
standards compliance, or public release preparation:

```
Audit this repo against pgEdge standards
```

```
Is this repo ready for release?
```

```
Bring this repo up to pgEdge standards
```

Or invoke directly:

```
/pgedge-skills:pgedge-repo-ready
```

## Re-run Behavior

The skill is idempotent. Re-running it on an already
onboarded repo detects drift — passing checks stay passing,
regressions are flagged, and only non-passing checks are
offered for remediation.

## File Structure

```
skills/pgedge-repo-ready/
  SKILL.md               # Skill orchestrator and phase routing
  README.md              # This file
  PLACEHOLDERS.md        # Placeholder reference
  standards/
    REPO_STANDARDS.md    # Canonical checklist (51 checks)
  phases/
    01-audit.md          # Detection and check execution
    02-report.md         # Terminal summary and file report
    03-remediation.md    # Template application and fixes
    04-claude-md-interview.md  # CLAUDE.md generation
  templates/
    ci/                  # GitHub Actions workflows
    claude/              # CLAUDE.local.md and agent defs
    docs/                # MkDocs Material scaffolding
    github/              # CODEOWNERS, issue/PR templates
    license/             # PostgreSQL License
    linting/             # golangci-lint, pre-commit, ruff
    review/              # CodeRabbit, Codacy configs
```
