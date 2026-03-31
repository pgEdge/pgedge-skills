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

## Support & Resources

For more information, visit
[docs.pgedge.com](https://docs.pgedge.com).

To report an issue, visit:
[github.com/pgEdge/pgedge-skills/issues](https://github.com/pgEdge/pgedge-skills/issues).

## License

This project is licensed under the
[PostgreSQL License](LICENSE.md).
