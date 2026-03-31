---
name: pgedge-docs
description: "Use this skill whenever the user wants to write, edit, review, or structure pgEdge documentation. This includes: creating new documentation files, editing existing docs, setting up MkDocs project structure from the pgEdge template, writing README files, formatting markdown content, creating changelogs, structuring navigation, or reviewing docs for style compliance. Trigger whenever the user mentions writing docs, documentation, README, mkdocs, changelog, or asks to review/fix doc formatting. Also trigger when setting up a new project's documentation scaffolding. Do NOT trigger for code-only changes that happen to be in markdown (e.g., editing a CLAUDE.md or DESIGN.md with development guidelines)."
---

# pgEdge Documentation Writing Skill

This skill ensures all pgEdge documentation follows consistent style,
formatting, and structural conventions. These rules apply across all
pgEdge projects and repositories.

## Project Setup from Template

When setting up documentation for a new pgEdge project, copy the
template files from this skill's `template/` directory (located at
`~/.claude/skills/pgedge-docs/template/`) into the target project.

### Template Contents

The template provides a minimal MkDocs Material configuration with
pgEdge branding:

```
template/
  mkdocs.yml                          # MkDocs Material config
  LICENSE.md                          # Root copy of PostgreSQL License
  docs/
    index.md                          # Landing page
    changelog.md                      # Keep a Changelog format
    LICENSE.md                        # PostgreSQL License (for MkDocs)
    img/
      favicon.ico                     # pgEdge favicon
      logo-dark.png                   # Dark mode logo
      logo-light.png                  # Light mode logo
    overrides/
      partials/
        logo.html                     # Light/dark logo switcher
    stylesheets/
      extra.css                       # Light/dark mode CSS
```

### Setting Up a New Project

When creating documentation scaffolding for a new project:

1. Copy the entire `template/` contents into the project root.
2. In `mkdocs.yml`, replace `<PROJECT NAME>` with the actual project
   name and `<PROJECT REPO>` with the GitHub repository name.
3. In `docs/index.md`, replace `<PROJECT NAME>` with the project name.
4. In `docs/changelog.md`, replace `<PROJECT NAME>` with the project
   name.
5. Build out the `nav:` section in `mkdocs.yml` to match the project's
   documentation structure.
6. Create the corresponding `.md` files in the `docs/` directory.

### mkdocs.yml Configuration

The template `mkdocs.yml` uses MkDocs Material with:

- `admonition`, `pymdownx.details`, and `pymdownx.superfences`
  markdown extensions.
- Light/dark mode palette switching with pgEdge branding.
- Custom logo partial for light/dark mode logos.
- Navigation features: instant loading, tracking, pruning, scroll to
  top, and TOC following.

The `nav:` section defines the left-hand table of contents. Keep it
synchronized with the actual files in `docs/`.

### Integration with the Primary Documentation Site

The primary pgEdge documentation site (the `pgedge-docs` repository,
published at docs.pgedge.com) uses the MkDocs `multirepo` plugin to
import documentation from individual project repositories. Each
project is imported via a `!import` directive in the primary site's
`mkdocs.yml` referencing the project's GitHub repo and branch:

```yaml
- v1.0: '!import https://github.com/pgEdge/project.git?branch=v1.0'
```

This means each project's `mkdocs.yml` and `docs/` folder must be
self-contained and valid on their own. Do not add configuration that
only works in the primary site context (e.g., `multirepo` plugin,
`redirects` plugin, analytics, consent banners, `nav_categories`,
or `versioned_docsets`).

## Product Names

The following product names are proper nouns; treat them as such. As a
rule, omit "the" in front of the name unless using it as an adjective
to describe software, files, or other artifacts:

- pgEdge Cloud
- pgEdge Enterprise Postgres
- pgEdge Distributed Postgres
- Enterprise Postgres
- Distributed Postgres
- pgEdge Postgres
- pgEdge AI Toolkit
- Spock
- LOLOR (aka lolor)
- ACE
- pgEdge Vectorizer
- Snowflake
- pgEdge Postgres MCP Server
- pgEdge Anonymizer
- pgEdge RAG Server
- pgEdge Docloader

The exception is "the Control Plane" which always takes the article.

## File Naming and Organization

### File Names

- Documentation files use the form `my_file.md` (lowercase,
  underscores).
- All docs live in the `docs/` folder or its subdirectories.
- Root-level markdown files use UPPERCASE names (`README.md`,
  `CONTRIBUTING.md`).
- The `LICENSE.md` file must be stored inside the `docs/` folder;
  aliases from other folders do not work properly with MkDocs.

### Directory Layout

Each project should have:

```
docs/
  index.md            # Main landing page
  changelog.md        # Release notes in Keep a Changelog format
  LICENSE.md          # PostgreSQL License
  img/                # Images and branding assets
  overrides/          # MkDocs Material overrides
  stylesheets/        # Custom CSS
```

## Line Length and Wrapping

- Wrap lines at **79 characters**.
- Lines should be as long as possible while staying within the limit.
- **Never** split hyperlinks across lines. The complete markdown link
  syntax `[text](URL)` must remain unbroken on a single line, even if
  it exceeds 79 characters. Break prose before or after the link.
- **Never** split table columns across lines.
- Replace em-dashes with regular dashes (hyphens).

## Writing Style

### Voice and Tone

- Write in **active voice**.
- Use full, grammatically correct sentences between 7 and 20 words.
- Use a semicolon to link similar ideas within a sentence.
- Use articles (a, an, the) when appropriate.
- Do not refer to an object as "it" unless the referent is in the same
  sentence. The pronoun is too ambiguous.

### What to Avoid

- Do not use bold font as headings; bold within inline text is fine.
- Do not use hashtags in front of steps in numbered lists.
- Do not use fragments in bulleted or stepped lists.
- Remove bold highlighted text content that is not inline; some
  software misinterprets it as a heading.

## Heading Structure

- Each file has **one** first-level heading (`#`).
- Use multiple second-level headings (`##`) for sections.
- Use third and fourth-level headings for sub-sections beneath the
  hierarchically previous heading level.
- Every heading must be followed by an introductory sentence or
  paragraph explaining what the section covers.
- **Never** use bold font within a heading.

## Bulleted Lists

- Use a lead-in sentence before the list that ends with a colon.
- The lead-in and bullet items together should form grammatically
  complete sentences.
- Keep lead-in sentences concise; use "include:" or "is useful for:"
  rather than verbose phrases like "include the following operations:".
- Items that complete the opening sentence start with a lowercase
  letter (unless a proper noun or code element) and end with a period.
- Items that are fragments or do not complete the introductory sentence
  start with uppercase and do not end in punctuation.
- Bullet items end only with periods, never semicolons.
- When bullet items are gerunds, they should work grammatically with
  the lead-in sentence.
- Always leave a blank line before the first item in any list.

Correct example:

```markdown
The view includes:

- the total entries and expired entries.
- the total cache size in megabytes.
- the average access count per entry.
```

Incorrect example:

```markdown
The view includes the following metrics:

- The view shows the total entries and expired entries.
- The view shows the total cache size in megabytes.
- The view shows the average access count.
```

## Numbered Lists

- Reserve numbered lists for sequential procedures where order is
  critical.
- Convert non-sequential content (feature lists, examples,
  comparisons) to bulleted lists or narrative format.
- Use an introductory sentence before each numbered list.
- Leave a blank line before the first item.
- Each entry should be a complete, correct sentence.
- Indent code blocks properly within numbered list items.
- Do not use bold font in numbered lists.

## Code Formatting

- Use an explanatory sentence before every code block that ends with a
  colon, identifying the command and describing what it does.
- Use single backticks for a single command or line of code:
  `` `SELECT * FROM my_table;` ``
- Use fenced code blocks with a language tag for multi-line code:

  ````markdown
  ```sql
  SELECT * FROM code;
  ```
  ````

- `stdio`, `stdin`, `stdout`, and `stderr` should always be in
  backticks.
- Capitalize SQL keywords; use lowercase for variables.
- Include links to third-party software documentation where referenced.
- Include links to the pgEdge GitHub repo when referring to cloning or
  working on the project.
- Do **not** create links to github.io.

## Tables

- Every table must have an introductory sentence immediately before it,
  using the pattern: "The following table [describes/shows/compares]
  X:".
- No icons or emojis in tables. Replace emoji characters with
  descriptive text (Fast, Good, High, etc.).
- Use plain text only; no bold font in cells.
- Avoid special characters that may not render in all formats.
- Use words rather than symbols to convey meaning.

## Hyperlinks

- The complete `[text](URL)` syntax must stay on one line.
- Break prose before or after the link, never within it.
- This rule supersedes the 79-character line length limit.

## Features and Overview Sections

If a page has a "Features" or "Overview" section after the introductory
paragraph, do not start with a heading. Instead, use a sentence like:
"The MCP Server includes the following features:", followed by a
bulleted list.

## Troubleshooting Sections

- Add troubleshooting content to a dedicated Troubleshooting section,
  not at the end of multiple individual doc files.
- Sort the Troubleshooting section and use sub-sections for topics
  like connection issues, authentication issues, and API-related
  issues.

## Next Steps Sections

- Include a "Next Steps" section (`## Next Steps`) at the end of
  tutorial and conceptual guide documents.
- Provide 2-4 bulleted links to related documentation.
- Each bullet should be a complete sentence describing what the linked
  document covers, using the pattern: "- The Document Name document
  describes/explains/provides X.".
- Do **not** include "Next Steps" in reference docs, API docs, FAQs,
  or troubleshooting guides.
- For FAQ and troubleshooting documents, use a support/contact section
  instead (e.g., "Still Have Questions?").

## README File Structure

Every public repository should have a `README.md` with the following
structure.

### Header

- A first-level heading (`#`) with the project name.
- Links to regression testing or developer tools if applicable.
- GitHub Action badges for important actions.

### Table of Contents

- A `## Table of Contents` heading.
- A bulleted list linking to files in `docs/`.
- Use gerunds in section titles where possible.

The following section titles are exceptions to the gerund rule:

- Architecture
- Overview
- Getting Started
- Release Notes
- FAQ
- Developer Resources
- Reference sections

### Suggested TOC Structure

Top-level entries should mirror the `mkdocs.yml` nav and link to the
same `.md` files; second-level links may reference anchors within a
file (`docs/file_name.md#section_name`):

- Architecture Guide
- Best Practices Guide
- Getting Started/Quick Start
- Building the Project (including prerequisites)
- Installing the Project
- Configuring the Project (with Advanced Configuration)
- Using the Project
- Upgrading the Project Installation
- Managing an Installation
- Modifying a Deployed Project
- Monitoring/Logging
- Performance
- Function Reference (with links to individual function pages)
- API Reference
- Troubleshooting
- FAQ
- Release Notes
- Developer Resources

### After the TOC

Include a short introduction describing the project, what it does, and
why someone might use it.

### Standard Sections

The following sections should appear in order:

**`## Installation`** - An introductory sentence, then:

- Prerequisite software or configuration.
- Installing with pgEdge software (e.g., the Control Plane) if
  applicable.
- Building from source.

**`## Configuration`** - An introductory sentence, then:

- Simple configuration steps.
- A link to advanced configuration options.

**`## Using <Project Name>`** - An introductory sentence, then:

- Simple usage instructions and examples.
- Links to advanced usage pages.

**`## Documentation`** - An introductory sentence, then:

- Details about building documentation from source.
- Links to the pgEdge documentation site.

**`## Support & Resources`**:

- Links to help and documentation: "For more information, visit here."
- A link to the Issues page: "To report an issue with the software,
  visit:" followed by the link.

**`## Contributing`** (if applicable):

- "We welcome your project contributions; for more information, see
  docs/developers.md."

Include a link to docs.pgedge.com: "For more information, visit
[docs.pgedge.com](https://docs.pgedge.com)."

**`## License`** (at the end):

- "This project is licensed under the
  [PostgreSQL License](LICENSE.md)."

## Copyright Notice

The copyright notice in `mkdocs.yml` uses the format
`Copyright &copy; 2025 - <current year> pgEdge, Inc`. When setting up
a new project or editing an existing one, always update the end year
to the current year. The template uses `<CURRENT YEAR>` as a
placeholder; replace it with the actual year at setup time.

## Changelog Format

Use the [Keep a Changelog](https://keepachangelog.com/) format with
[Semantic Versioning](https://semver.org/spec/v2.0.0.html). The
template file at `template/docs/changelog.md` provides the starting
structure.

## mkdocs.yml Navigation Structure

The `nav:` section in `mkdocs.yml` defines the left-hand navigation
pane. It should mirror the README TOC and may include:

```yaml
nav:
  - Architecture Guide
  - Best Practices Guide
  - Installing this Project:
      - Getting Started/Quick Start
      - Building the Project:
          - Prerequisites
      - Installation
      - Configuring this Project:
          - Advanced Configuration
      - Upgrading the Project
  - Using the Project:
      - Connecting to the Project
      - Project Usage Instructions
  - Using Project Management Features:
      - Modifying the Installation
      - Monitoring/Logging
      - Performance
  - Function Reference
  - API Reference
  - Troubleshooting
  - FAQ
  - Release Notes
  - Developer Resources
```

Nested details can be in the same file as the parent section (e.g.,
Prerequisites can live in the same file as Installation).

## Review Checklist

When reviewing documentation, verify:

- [ ] One `#` heading per file
- [ ] Every heading followed by an introductory sentence
- [ ] Lines wrapped at 79 characters (hyperlinks excepted)
- [ ] No bold used as headings
- [ ] No em-dashes (use hyphens)
- [ ] No fragments in lists
- [ ] Bulleted lists have lead-in sentences ending with colons
- [ ] Numbered lists used only for sequential procedures
- [ ] Code blocks have language tags and explanatory sentences
- [ ] Tables have introductory sentences and no emojis
- [ ] Product names treated as proper nouns
- [ ] No ambiguous "it" references
- [ ] Hyperlinks not split across lines
- [ ] Blank lines before all lists
- [ ] Active voice throughout
