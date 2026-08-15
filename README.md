# pgEdge Skills

A Claude Code plugin providing pgEdge-specific skills for
documentation writing, and more to come.

## Installation

Install the plugin with the following commands in Claude Code:

```
/plugin marketplace add pgEdge/pgedge-skills
/plugin install pgedge-skills@pgedge-skills
```

## Skills

### pgedge-docs

Documentation writing skill that:

- enforces pgEdge documentation style and formatting conventions.
- scaffolds new MkDocs Material projects with pgEdge branding.
- guides README structure, changelog format, and navigation layout.
- validates documentation against a comprehensive review checklist.

The skill activates automatically when you work on documentation
tasks. You can also invoke it explicitly:

```
/pgedge-skills:pgedge-docs
```

#### Example Prompts

**Scaffolding a new docs site:**

```
Set up MkDocs documentation scaffolding for the pgEdge Vectorizer
project in ~/git/vectorizer. The docs should cover installation,
configuration, usage, and a function reference.
```

**Reviewing existing docs:**

```
Review the docs/ directory for style compliance. Check line length,
heading structure, list formatting, and product name usage, and fix
any issues you find.
```

**Documenting a new feature:**

```
Document the new cache_stats view for Spock. It exposes total
entries, expired entries, cache size in MB, and average access
count. Add it to the existing monitoring page and update the
changelog.
```

### [pgedge-repo-ready](skills/pgedge-repo-ready/README.md)

Repository audit and remediation skill that:

- audits repos against 51 pgEdge engineering standards across
  6 tiers (governance, code quality, docs, developer experience,
  language-specific, Postgres).
- detects languages (Go, Python, TypeScript), project type, and
  Postgres usage to determine which checks apply.
- applies 35 bundled templates and generates detection-driven
  fixes for CI, linting, documentation, and devcontainer setup.
- generates CLAUDE.md via structured interview if missing.
- is idempotent — re-run to detect drift or regressions.

The skill activates when you mention repo readiness, standards,
or release preparation. Invoke explicitly with:

```
/pgedge-skills:pgedge-repo-ready
```

#### Example Prompts

**Full audit and remediation:**

```
Audit this repo against pgEdge standards and fix everything.
```

**Audit only (no changes):**

```
Is this repo ready for public release?
```

**Targeted fix:**

```
Bring just the CI and governance up to pgEdge standards.
```

### pgedge-psql (beta)

PostgreSQL `psql` skill that:

- guides efficient use of the `psql` command-line tool.
- defaults to compact, token-efficient output formats.
- enforces read-only access patterns unless writes are needed.
- provides schema discovery and connection workflows.

The skill activates automatically when working with PostgreSQL
databases via `psql`. You can also invoke it explicitly:

```
/pgedge-skills:pgedge-psql
```

### pgedge-webapp (beta)

Full-stack web app scaffolding skill that:

- generates a Go (sqlite, bcrypt) backend plus React (MUI, TypeScript)
  frontend.
- locks in pgEdge design conventions: theme, header, login, form-label
  CSS, accessibility patterns.
- ships a Helm chart, Dockerfile, docker-compose, and CI workflows.
- enforces 90% test coverage via vitest thresholds and a Go coverage
  gate.

The skill activates automatically when starting a new full-stack
project. Invoke explicitly with:

```
/pgedge-skills:pgedge-webapp
```

#### Example Prompt

```
Scaffold a new pgEdge web app at ~/git/my-app called "My App".
```

### pgedge-compactor

Skill-file compactor that:

- classifies content blocks as routing, priming, examples, tables,
  cross-routing, or frontmatter.
- compresses only routing blocks while preserving all
  quality-critical content.
- converts parallel prose into tables where applicable, including
  the legend technique for repeated values.
- when prose compression is exhausted, structurally extracts large
  or conditional content into files that load only when needed.
- supports single file, directory (two-pass resolution), and glob
  targets.

The skill activates automatically when you ask to compact or
optimize a skill file. You can also invoke it explicitly:

```
/pgedge-skills:pgedge-compactor
```

#### Example Prompts

**Compacting a single skill file:**

```
Compact the pgedge-docs SKILL.md to reduce token usage.
```

**Compacting all skills in a directory:**

```
Compact all skill files in the skills/ directory.
```

**Reviewing compaction potential:**

```
How much could we compress the pgedge-webapp skill? Show me
what's routing vs priming.
```

### pgedge-plain-speak

Plain-language skill that:

- replaces jargon used as ordinary vocabulary with everyday words.
- keeps technical terms that are the actual name for something, such
  as race condition, deadlock, or memoization.
- prefers the shorter everyday word wherever it means the same thing.
- applies to code comments as well as prose.
- changes the wording without changing the depth of an explanation.

The skill activates when you ask for ordinary English instead of
jargon-flavored phrasing. You can also invoke it explicitly:

```
/pgedge-skills:pgedge-plain-speak
```

#### Example Prompts

**Asking for a plain explanation:**

```
Explain why this query is slow, in plain English.
```

**Cleaning up existing prose:**

```
Rewrite the comments in this file without the jargon.
```

## Support & Resources

For more information, visit
[docs.pgedge.com](https://docs.pgedge.com).

To report an issue, visit:
[github.com/pgEdge/pgedge-skills/issues](https://github.com/pgEdge/pgedge-skills/issues).

## License

This project is licensed under the
[PostgreSQL License](LICENSE.md).
