# Changelog

All notable changes to pgedge-skills are documented here. The format is
based on [Keep a Changelog](https://keepachangelog.com/), and this
project adheres to semantic versioning.

## [1.2.5] - 2026-07-02

### Added

- Automated release pipeline: a push to `main` that bumps the plugin
  version builds per-skill zips and publishes a GitHub Release, gated on
  lint (`.github/workflows/ci.yml`, `scripts/build-skill-zips.sh`).
- Layered markdownlint config: a permissive skill baseline plus a strict
  pgEdge house-style override for the pgedge-docs exemplar docs.
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
