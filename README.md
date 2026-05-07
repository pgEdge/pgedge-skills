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

## Support & Resources

For more information, visit
[docs.pgedge.com](https://docs.pgedge.com).

To report an issue, visit:
[github.com/pgEdge/pgedge-skills/issues](https://github.com/pgEdge/pgedge-skills/issues).

## License

This project is licensed under the
[PostgreSQL License](LICENSE.md).
