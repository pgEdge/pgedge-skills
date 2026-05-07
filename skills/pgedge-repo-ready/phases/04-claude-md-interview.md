# Phase 4 — CLAUDE.md Interview

Generate a structured CLAUDE.md and set up the `.claude/`
directory through a guided Q&A chain. This phase runs when
CLAUDE.md is MISSING or PARTIAL after remediation.

## Pre-Interview Detection

Before asking questions, gather what can be detected
automatically. This seeds the interview with informed defaults.

```bash
# Project name (already resolved in Phase 3)
# Languages (already detected in Phase 1)
# Makefile targets
grep "^[a-z].*:" Makefile 2>/dev/null | cut -d: -f1

# Directory structure (top 2 levels)
find . -maxdepth 2 -type d \
  -not -path "./.git*" \
  -not -path "./.claude*" \
  -not -path "./node_modules*" \
  -not -path "./vendor*"

# Test framework
grep -r "testing" --include="*.go" -l 2>/dev/null | head -3
grep -r "pytest\|unittest" --include="*.py" -l 2>/dev/null | head -3
grep -r "vitest\|jest" package.json 2>/dev/null

# Entry points
grep -r "func main()" --include="*.go" -l 2>/dev/null
grep "entry_points" pyproject.toml 2>/dev/null

# Docker services
grep "^\s\+[a-z].*:" docker-compose.yml 2>/dev/null \
  | grep -v "^\s*#"

# Existing CLAUDE.md sections (if PARTIAL)
grep "^##" CLAUDE.md 2>/dev/null
```

## Interview Questions

Ask one question at a time. Wait for the user's answer before
proceeding to the next question. Offer multiple choice where
possible, seeded from detection results.

### Question 1: Project Purpose

```
What is this project?

Based on what I detected:
  Languages: <detected>
  Structure: <detected dirs>

  A) CLI tool / utility
  B) Library / package
  C) Web service / API
  D) Multi-component system (multiple services)
  E) Something else (describe)
```

Record answer as `PROJECT_TYPE_DESCRIPTION`.

### Question 2: Architecture

```
Describe the architecture in 2-3 sentences.

I detected these components:
  <list detected directories with go.mod, package.json, or
   Dockerfile>
  <list docker-compose services if present>

How do these relate to each other? What's the data flow?
```

Record answer as `ARCHITECTURE_DESCRIPTION`.

### Question 3: Key Commands

```
What are the key development commands?

I found these Makefile targets:
  <list detected targets>

Are there additional commands developers need to know? Any
setup steps, database migrations, or environment requirements?
```

Record answer as `KEY_COMMANDS`.

### Question 4: Testing Strategy

```
What's the testing strategy?

I detected:
  Test framework: <detected>
  Test directories: <detected>
  Coverage tool: <detected or "none found">

A few questions:
  - Are there integration tests that need external services
    (database, Docker)?
  - What's the coverage expectation?
  - Any test commands beyond what's in the Makefile?
```

Record answer as `TESTING_STRATEGY`.

### Question 5: Primary Maintainers

```
Who are the primary maintainers?

This is used for the CODEOWNERS file cross-reference and
the CLAUDE.md context section. List GitHub usernames and
areas of ownership.

Example: @username1 (backend), @username2 (frontend)
```

Record answer as `MAINTAINERS`.

### Question 6: Sub-Agents

```
Which Claude sub-agents would be useful for this repo?

Based on detected languages and patterns, I'd suggest:
  <if Go detected>
  [x] golang-expert — Go development, testing, MCP protocol
  <if Python detected>
  [x] python-expert — Python development, testing
  <if TypeScript detected>
  [x] react-expert — React/TypeScript, component architecture
  <if Postgres detected>
  [x] postgres-expert — Database admin, Spock replication
  <always>
  [x] documentation-writer — Markdown docs, style compliance
  [x] security-auditor — Security review, vulnerability detection

Which of these should I include? Add or remove as needed.
```

Record answer as `SELECTED_AGENTS`.

### Question 7: Project-Specific Conventions

```
Any project-specific conventions not covered above?

Examples:
  - Specific API design patterns
  - Error handling conventions
  - Naming conventions beyond pgEdge defaults
  - Deployment requirements
  - Environment variable patterns
```

Record answer as `CUSTOM_CONVENTIONS`.

## CLAUDE.md Generation

Using the interview answers, generate CLAUDE.md with this
structure:

```markdown
# <PROJECT_NAME> — Instructions for Claude

## Context

<PROJECT_TYPE_DESCRIPTION — expanded into 2-3 sentences>

## Architecture

<ARCHITECTURE_DESCRIPTION>

<If multi-component, list components with one-line descriptions>

## Development

### Prerequisites

<Detected from languages, tools, docker-compose>

### Key Commands

<KEY_COMMANDS, formatted as a table or list>

### Testing

<TESTING_STRATEGY>

## Code Standards

<Language-specific standards from REPO_STANDARDS.md>
<CUSTOM_CONVENTIONS if provided>

### Git Practices

- Conventional commit style: feat:, fix:, docs:, chore:, etc.
- Header <= 50 characters, body wrapped at 72 characters
```

Sections that could not be fully populated from the interview
should include a `<!-- REVIEW: description of what to add -->`
comment so the user can fill them in later.

## .claude/ Directory Setup

### Create Agent Files

For each agent in SELECTED_AGENTS, copy the corresponding
template from `templates/claude/agents/` to `.claude/agents/`.

Customize each agent file with project-specific context:
- Replace generic descriptions with project-relevant details
- Add any detected knowledge base files or patterns

### Create Directory Structure

```bash
mkdir -p .claude/agents .claude/specs .claude/plans
```

### Create CLAUDE.local.md

Copy `templates/claude/CLAUDE.local.md` to the repo root. Ensure
it is in `.gitignore`:

```bash
grep -q "CLAUDE.local.md" .gitignore 2>/dev/null \
  || echo "CLAUDE.local.md" >> .gitignore
```

## Completion

After CLAUDE.md and .claude/ are set up, print:

```
CLAUDE.md generated and .claude/ directory configured.

Files created:
  - CLAUDE.md (review <!-- REVIEW --> comments)
  - .claude/agents/<list created agents>
  - .claude/specs/ (empty, ready for use)
  - .claude/plans/ (empty, ready for use)
  - CLAUDE.local.md (gitignored)

Please review CLAUDE.md and fill in any <!-- REVIEW --> sections
with project-specific details.
```

Update REPO_READY_AUDIT.md with final status for X-01 through
X-05.
