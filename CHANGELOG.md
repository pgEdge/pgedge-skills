# Changelog

All notable changes to pgedge-skills are documented here. The format is
based on [Keep a Changelog](https://keepachangelog.com/), and this
project adheres to semantic versioning.

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
