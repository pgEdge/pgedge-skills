# Template Placeholders

Templates in the `templates/` directory use angle-bracket
placeholders that are substituted during remediation.

## Content Placeholders

Substituted in file contents via find-and-replace:

| Placeholder | Source | Example |
|---|---|---|
| `<PROJECT_NAME>` | go.mod module name, pyproject.toml name, or package.json name | `mm-ready` |
| `<GITHUB_REPO>` | git remote origin URL, extracted as `org/repo` | `pgEdge/mm-ready` |
| `<GITHUB_ORG>` | Organization from remote URL | `pgEdge` |
| `<CURRENT_YEAR>` | Current four-digit year | `2026` |
| `<PRIMARY_OWNER>` | User input — GitHub username of primary maintainer | `@AntTheLimey` |
| `<PRIMARY_LANGUAGE>` | Detected primary language | `Go` |
| `<PROJECT_DESCRIPTION>` | README first paragraph or user input | `Multi-master readiness scanner` |
| `<GO_VERSION>` | From go.mod `go` directive, or latest stable | `1.26` |
| `<NODE_VERSION>` | From .nvmrc, package.json engines, or latest LTS | `22` |
| `<PYTHON_VERSION>` | From pyproject.toml requires-python lower bound | `3.10` |
| `<PYTHON_VERSIONS>` | Matrix of supported Python versions | `["3.10", "3.11", "3.12", "3.13"]` |
| `<PG_VERSIONS>` | PostgreSQL version matrix for CI | `["16", "17", "18"]` |
| `<COPYRIGHT_RANGE>` | Year range for license | `2025-2026` |

## Detection Priority

When auto-detecting `<PROJECT_NAME>`:

1. `go.mod` — module path, last segment
2. `pyproject.toml` — `[project] name`
3. `package.json` — `name` field
4. Git remote — repo name from origin URL
5. Directory name — fallback

When auto-detecting `<GITHUB_REPO>`:

1. `git remote get-url origin` — parse org/repo from URL
2. User input — if no remote configured

## Notes

- All placeholders use `<UPPER_SNAKE_CASE>` format
- Substitution is a simple string find-and-replace across all
  files in the copied template subtree
- Verify no unsubstituted `<` + `UPPER_SNAKE` patterns remain
  after substitution
- Template files that are language-specific (e.g., ci-go.yml) are
  only copied when that language is detected
