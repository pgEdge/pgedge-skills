---
name: pgedge-repo-ready
description: "Use when the user wants to audit a repository against pgEdge engineering standards, prepare a repo for public release, check repo compliance, onboard a new repo to pgEdge conventions, or detect drift in a previously onboarded repo. Trigger on phrases like 'repo ready', 'audit this repo', 'prepare for release', 'check standards', 'is this repo ready', 'bring up to standard', 'onboard this repo'. Also trigger when the user mentions CODEOWNERS, CI workflows, CodeRabbit, pre-commit hooks, or coverage thresholds in the context of repo setup. Do NOT trigger for code changes within an already-conforming repo, for adding a single CI job, or for documentation-only edits (use pgedge-docs for those)."
---

# pgEdge Repo Ready

Audit a repository against pgEdge engineering standards and guide
remediation to bring it to release quality. Re-runnable for drift
detection.

## How This Skill Works

This skill has four phases executed in sequence. Each phase is
defined in a separate document under `phases/`. The canonical
standard is in `standards/REPO_STANDARDS.md`. Templates for
remediation are in `templates/`.

**Do not embed standards, audit logic, or template content in this
file.** Route to the appropriate phase document instead.

## Phase Sequence

### Phase 1 — Audit

Read `phases/01-audit.md` and follow its instructions.

Scan the current working directory against
`standards/REPO_STANDARDS.md`. Detect languages, project type,
and Postgres interaction. Run every applicable check. Record each
result as PASS, PARTIAL, MISSING, or SKIP.

Do not create, modify, or delete any files during this phase.

### Phase 2 — Report

Read `phases/02-report.md` and follow its instructions.

Print a terminal summary with tier-by-tier scores. Write the
detailed report to `REPO_READY_AUDIT.md` in the repo root (add
it to `.gitignore` if not already present). Prompt the user to
choose a remediation mode.

### Phase 3 — Remediation

Read `phases/03-remediation.md` and follow its instructions.

Walk through each failing or partial check. For each one, explain
what will change, show the content, apply on confirmation, and
verify. Use templates from the `templates/` directory, applying
placeholder substitutions per `PLACEHOLDERS.md`.

### Phase 4 — CLAUDE.md Interview

Read `phases/04-claude-md-interview.md` and follow its
instructions.

If CLAUDE.md is MISSING or PARTIAL after remediation, run a Q&A
chain to generate a structured CLAUDE.md with architecture,
commands, testing strategy, and sub-agent definitions. Set up the
full `.claude/` directory (agents, specs, plans).

## Re-run Behavior

This skill is idempotent. When run on a previously onboarded repo:

- The audit phase checks current state against current standards
- Passing checks are reported normally (no remediation offered)
- Regressions are highlighted as REGRESSION
- New checks (from updated standards) are reported as NEW
- Remediation only covers non-passing checks
- REPO_READY_AUDIT.md is overwritten with current state
