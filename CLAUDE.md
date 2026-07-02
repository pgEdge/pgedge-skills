# pgedge-skills — Instructions for Claude

Context for Claude sessions working in this repository.

## Context

This repo is a Claude Code plugin marketplace of pgEdge skills. It ships
a small collection of task-focused skills that bring pgEdge conventions
and tooling into a Claude Code session:

- `pgedge-docs` — write, edit, and structure pgEdge documentation
  (MkDocs setup, README/changelog authoring, style compliance).
- `pgedge-psql` (beta) — query and explore a PostgreSQL database with
  the `psql` CLI.
- `pgedge-webapp` (beta) — scaffold and build a pgEdge web application
  from a shipped template.

`pgedge-compactor` arrives later via rebase and is not yet present on
this branch — do not be surprised when it appears.

Users install the marketplace with
`/plugin marketplace add pgEdge/pgedge-skills` and receive updates via
`/plugin update`.

## Architecture

- `skills/` holds one directory per skill. Each skill has a `SKILL.md`
  routing layer — a concise front door that tells Claude when to
  activate and where to look. Detailed content lives in `references/`
  under the skill.
- Some skills ship a `template/` scaffold that is copied into a
  destination project (e.g. `pgedge-docs` and `pgedge-webapp`). A skill
  may also carry `tools/` for helper scripts (e.g. `pgedge-webapp`).
  Template payload under `skills/**/template/**` follows the destination
  project's conventions, not this repo's — it is excluded from lint at
  the CI-glob level.
- `.claude-plugin/` holds the plugin metadata: `plugin.json` (name,
  description, `version`, author) and `marketplace.json` (the plugin
  listing consumed by `/plugin marketplace add`).

## Development

### Skill edits

- **Use Opus** for writing new skill prose or designing a skill's
  routing layer. **Sonnet is fine** for mechanical edits — replacements,
  reference wiring, structural moves — when the exact content is
  specified in a plan.
- `SKILL.md` is the routing layer. Keep it concise: it is loaded to
  decide activation, so every token there is paid often. Push detailed
  content down into `references/`.
- **Inline vs extract:** every file read costs roughly 90 tokens of
  fixed overhead (tool call + wrapper + line numbers). Do not extract
  short, always-needed content into a separate file — that makes
  activation slower and more expensive. Extract only content that is
  large or conditional (needed in a minority of activations).

### Testing

- When running benchmarks via subagents, each subagent writes its own
  result file to disk as its final action, so the result survives even
  if the parent never reads the return value.
- **Token counts alone do not prove quality.** A cheaper run that
  produces worse output is not an improvement — preserve full response
  text and compare it against the baseline, not just cost.
- Prefer A/B comparisons over absolute scores. `pgedge-compactor` was
  validated by a blind A/B benchmark rather than by a single-number
  target.

### Releasing

This repo is a Claude Code plugin marketplace. Users install it with
`/plugin marketplace add pgEdge/pgedge-skills` and receive updates via
`/plugin update`.

#### The version bump is the update signal — nothing else is

Claude Code fetches plugin content from the **default branch (`main`)
HEAD**. It decides whether an update is available by comparing the
`version` field in `.claude-plugin/plugin.json` against what the user
has installed. Consequences:

- **Bump `version` in `plugin.json` for every user-facing change.**
  Pushing commits to `main` without bumping it does nothing — users
  see "already at the latest version."
- Git tags and GitHub Releases play **no part** in distribution. The
  `/plugin` system never downloads release assets; it clones the repo
  at HEAD. Releases and the per-skill zips are informational only.
- `version` lives in `plugin.json` **only**. Do not add a `version` to
  `marketplace.json` — if the two drift, Claude Code silently uses the
  `plugin.json` value and the manifest version is ignored, which is a
  hard-to-spot bug.

#### Release flow

1. Make your change on a branch; bump `version` in `plugin.json`
   (semver: patch by default).
2. Add a matching `## [x.y.z]` entry to `CHANGELOG.md`.
3. Open a PR to `main`. CI lint must pass.
4. On merge to `main`, the release workflow reads the version, and if
   no `vX.Y.Z` tag exists yet, builds the per-skill zips and publishes
   a GitHub Release with auto-generated notes.

The merge that bumps the version is therefore a single event that both
signals the update to users and cuts the release — they can't drift
apart.

## Code Standards

### Skill-file formatting

Skill-authored content — `SKILL.md`, everything under `references/**`,
and any `knowledge-base.md` — is **never hard-wrapped**. Write one line
per paragraph and let your editor soft-wrap. Every hard line break is a
newline token that Claude pays on **every activation** of the skill, so
wrapping this content makes it permanently more expensive to run.

This is a **convention, not a CI gate.** markdownlint here does not
police line width anywhere (MD013 is off), and there is no wrap-check in
CI — nothing will fail if you wrap a line. The 79-column pgEdge doc
style applies to *documentation*, where it is enforced by the
`pgedge-docs` skill and downstream CI, not to skill prose in this repo.

### Commits

- Use conventional commit style: `feat:`, `fix:`, `docs:`, `chore:`,
  `ci:`.
- Prefer small, focused commits over large batched changes.

### Content provenance

Skills in this repo embed first-party pgEdge documentation conventions
and original content. Keep them first-party — do not fold in third-party
material or content of unclear origin.

### Git Practices

- **Never** add self-attribution to commits or PRs — no `Co-Authored-By`
  trailers, and no mention of Claude, AI, or LLM tooling in commit
  messages, PR titles, or PR bodies.
- Keep commits small and focused, each doing one thing.
