# Changelog

All notable changes to pgedge-skills are documented here. The format is
based on [Keep a Changelog](https://keepachangelog.com/), and this
project adheres to semantic versioning.

## [1.2.8] - 2026-08-20

### Changed

- `pgedge-repo-ready` G-06 now decides whether a repository needs its
  own `SECURITY.md` from the org custom property `repo_type`, rather
  than from a reading of what the repository appears to be. Only
  `repo_type=product` requires an in-repo copy; every other value
  passes on the organisation default in `pgEdge/.github`. The audit
  reads the property instead of inferring it.
- A repository with no `repo_type` value now reports `UNCLASSIFIED`,
  which is neither a pass nor a fail. Most of the organisation carries
  no value yet, so absence is not evidence that a repository is not a
  product. The audit flags it and asks the owner to set the property.

## [1.2.7] - 2026-08-04

### Changed

- `pgedge-repo-ready` G-06 now issues the approved pgEdge security
  policy. The shipped `SECURITY.md` template previously promised
  acknowledgement within 48 hours and an estimated timeline for a fix,
  neither of which is approved, so every onboarded repository published
  a commitment pgEdge had not made. It now points reporters at
  `security@pgedge.com` and the Vulnerability Disclosure Statement at
  `docs.pgedge.com/security`, and states the approved five-business-day
  acknowledgement.
- G-06's location moved from `.github/SECURITY.md` to the repository
  root, which is where tooling that checks for a security policy looks.
  A repository that is not a product now passes G-06 on the
  organisation-wide default in `pgEdge/.github` and needs no file of its
  own; a pre-existing `.github/SECURITY.md` is reported as being in the
  wrong location.

## [1.2.6] - 2026-07-07

### Added

- `pgedge-compactor` skill: compacts Claude Code skill files to
  reduce token cost while preserving priming content, concrete
  examples, tables, and cross-routing prompts. Compresses routing
  prose (Axis 1) and, when prose is exhausted or a file mixes
  always-needed with conditional content, structurally extracts
  content into files that load only when needed (Axis 2). Ships a
  knowledge base with the compression taxonomy, tabularization
  guide, and structural-extraction method.

## [1.2.5] - 2026-07-02

### Added

- Automated release pipeline: a push to `main` that bumps the plugin
  version builds per-skill zips and publishes a GitHub Release, gated on
  lint (`.github/workflows/ci.yml`, `scripts/build-skill-zips.sh`).
- Permissive markdownlint baseline (`.markdownlint.yaml`) that checks
  well-formedness, not line width; template payload
  (`skills/**/template/**`) is excluded from lint.
- Root `CLAUDE.md` documenting repo context and the release flow.

### Changed

- Version is now single-sourced in `plugin.json`; removed the `version`
  field from `marketplace.json`.

## [1.2.4] - 2026-05-05

- User-management endpoints in the pgedge-webapp template.

## [1.2.1] - 2026-05-05

- Docs-layout fix (`docs/api/*` instead of `docs/docs/api/*`) and
  `make-openapi` target / `browser.md` frontmatter fixes from e2e
  testing.
