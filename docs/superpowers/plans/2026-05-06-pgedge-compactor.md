# pgedge-compactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use
> superpowers:subagent-driven-development (recommended) or
> superpowers:executing-plans to implement this plan task-by-task.
> Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a skill-file compactor to the pgedge-skills plugin
that compresses routing prose while preserving priming content,
examples, and cross-routing prompts.

**Architecture:** Two markdown files — `SKILL.md` (process,
classification, rules) and `knowledge-base.md` (taxonomy,
patterns, anti-patterns) — in `skills/pgedge-compactor/`. The
skill loads the knowledge base on activation, resolves targets
(file, directory two-pass, or glob), then runs a per-file
classify-compress-verify-report loop. No code — pure skill
authoring.

**Tech Stack:** Claude Code skill plugin (markdown with YAML
frontmatter), JSON config files.

**Spec:** `docs/superpowers/specs/2026-05-06-pgedge-compactor-design.md`

**Source material:** `/tmp/skill-compactor-extracted/skill-compactor/`
contains the gm-apprentice originals (`agent.md`, `knowledge-base.md`)
to adapt from. All TTRPG references must be replaced with
coding/pgedge examples.

---

## File Map

| Action | File | Responsibility |
|--------|------|---------------|
| Create | `skills/pgedge-compactor/SKILL.md` | Process, classification taxonomy, compression techniques, rules, invocation, reporting |
| Create | `skills/pgedge-compactor/knowledge-base.md` | Compression taxonomy, safe/preserve lists, tabularization guide, diagnostic signals, anti-patterns |
| Modify | `README.md` | Add pgedge-compactor section following existing skill entry pattern |
| Modify | `.claude-plugin/plugin.json` | Bump version 1.2.4 → 1.3.0 |
| Modify | `.claude-plugin/marketplace.json` | Bump version 1.2.4 → 1.3.0 |

---

### Task 1: Create the knowledge-base.md

The knowledge base must exist before the SKILL.md because the
skill references it. Write it first so the cross-routing from
SKILL.md is valid.

**Files:**
- Create: `skills/pgedge-compactor/knowledge-base.md`

- [ ] **Step 1: Create the directory**

```bash
mkdir -p skills/pgedge-compactor
```

- [ ] **Step 2: Write knowledge-base.md**

Create `skills/pgedge-compactor/knowledge-base.md` with the
following content. This adapts the gm-apprentice knowledge base
with all TTRPG references replaced by coding/pgedge examples,
and adds the tabularization section from the spec.

```markdown
# Skill Compactor Knowledge Base

Compression patterns learned from skill compaction work.
Manually updated when new patterns are discovered.

## The Taxonomy

1. **Routing compresses** - text that tells the agent where to
   look. "Read knowledge-base.md in this directory" is routing.
   Connective tissue, introductory prose, "see X for Y",
   section transitions. All compressible.
2. **Priming preserves** - text that shapes output quality.
   "Enforce 90% coverage gate. Use `--csv -t` for queries."
   Named patterns, thresholds, specific numbers, quality
   criteria, checklist items. Preserve always.
3. **Concrete examples are sacred** - code samples, query
   examples, and usage demonstrations teach the agent what
   correct output looks like. Removing the psql format
   comparison table and flag examples caused the agent to
   default to verbose aligned output, wasting tokens. The
   examples are load-bearing - they calibrate the agent's
   sense of what "correct" means.
4. **Deduplication needs breadcrumbs** - when content lives in
   one canonical location, every file that used to have it
   needs a redirect, not silence. Agents that find nothing
   stop looking.

## What Compresses Safely

- Introductory paragraphs ("This skill guides you through...").
- "See X for Y" prose (shorten, don't remove the pointer).
- Section transitions and connective tissue.
- Multi-sentence routing that can become single-line arrows.
- Redundant phrasing saying the same thing twice.
- Verbose bullet explanations where the bullet text is clear.
- Multi-paragraph command descriptions reduced to single-line
  directives with compact priming.

## What Must Be Preserved

- **Named patterns** - Form Field Pattern, Read-Only Access
  Pattern, Token-Efficient Output Formats, Keep a Changelog
  format, Safe Write Pattern, or any named pattern in the
  target file. The name itself is a routing anchor.
- **Specific thresholds** - 90% coverage gate, 79-character
  line wrapping, port range 1024-65535, LIMIT 5 for initial
  exploration, 12-character minimum password. Any number that
  serves as a quality gate or constraint.
- **Tables** - already near-optimal density. Compress prose
  around tables, not the tables themselves.
- **Checklist items and coverage criteria** - quality gates the
  agent checks against.
- **Trigger phrases** - verbatim, these are routing anchors
  that determine when the skill activates.
- **Frontmatter** - never touch YAML frontmatter between `---`
  fences.
- **Concrete examples** - code samples, SQL queries, config
  snippets, bash commands, any content that demonstrates what
  correct output looks like for a specific domain.
- **Cross-routing prompts** - "For X, also read Y.md" - the
  biggest quality lever found in benchmarking. Compacted files
  feel "complete enough" so agents stop; cross-routing nudges
  them to keep going.

## Tabularization

Converting parallel prose into tables is one of the
highest-yield compression techniques. A well-formed table
preserves all information at a fraction of the token cost.

**Candidates for tabularization (3+ items required):**

- Lists of options, flags, or settings where each item has the
  same attributes (name, default, description, type).
- Comparison lists ("X does this, Y does that") where items
  share a consistent structure.
- Conditional behavior ("when A, do B; when C, do D") with
  consistent condition/action pairs.
- Format or mode descriptions with consistent fields. The psql
  skill's format summary table (Scenario | Flags | Why) is a
  model of good tabularization.
- Error code or status code mappings.
- Permission or role matrices.

**Do NOT tabularize:**

- Sequential procedures (numbered steps depend on order and
  context from previous steps - tables lose that).
- Code blocks and concrete examples (sacred content - tables
  strip formatting and context).
- Narrative explanations where meaning builds across sentences.
- Items with wildly different structures or field lengths (a
  table with one cell containing a paragraph defeats the
  purpose).
- Fewer than 3 items (a two-row table has more overhead than
  two bullets).

**How to structure the table:**

- Columns are the shared attributes. Rows are the items.
- Column headers should be single words or very short phrases.
- Cell content should be terse - if a cell needs more than ~10
  words, the content may not be a good table candidate.
- Preserve all information from the original prose. If
  tabularizing forces you to drop nuance, keep the prose.
- Add a one-line introductory sentence before the table.

**Legend technique for repeated values:**

When the same long values appear across many rows, replace them
with short codes and add a legend above the table.

Example - a permissions matrix before compression:

| Endpoint | superuser | administrator | editor | viewer |
|---|---|---|---|---|
| GET /users | read-write | read-write | read-only | read-only |
| POST /users | read-write | read-write | none | none |
| DELETE /users | read-write | none | none | none |
| GET /clusters | read-write | read-write | read-only | read-only |

After applying the legend technique:

**su**: superuser, **adm**: administrator, **ed**: editor,
**vi**: viewer - **RW**: read-write, **RO**: read-only,
**--**: none

| Endpoint | su | adm | ed | vi |
|---|---|---|---|---|
| GET /users | RW | RW | RO | RO |
| POST /users | RW | RW | -- | -- |
| DELETE /users | RW | -- | -- | -- |
| GET /clusters | RW | RW | RO | RO |

The legend costs a few tokens once; the savings compound across
every row. Apply when 3+ values repeat across 4+ rows.

## Diagnostic Signals

- **File resists compaction** - already well-structured. The
  psql skill's format summary table and schema discovery SQL
  blocks are 80% structured content. Minimal savings is
  success, not failure.
- **Verbose prose dilutes signal** - after compaction, actionable
  content is a larger proportion of what the agent reads. The
  model finds useful bits faster.
- **Tables compress beautifully** - they preserve information
  density at lower token cost than equivalent prose.

## Anti-Patterns

- Cutting priming to hit a size target (the routing/priming
  boundary is sacred - never cross it for space).
- Removing concrete examples because they seem like "filler"
  (output quality drops when the agent loses calibration
  examples).
- Deduplicating without leaving redirects (agents find nothing
  and stop looking).
- Compressing domain-specific terminology into generic terms
  (replacing "psql" with "database CLI" caused the agent to
  suggest alternative tools - domain terms are routing anchors).
- Removing cross-routing prompts as redundant (they're the
  single biggest quality lever found).
```

- [ ] **Step 3: Verify the file exists and has no YAML frontmatter**

The knowledge base is a reference file, not a skill entry point.
It should NOT have `---` frontmatter.

```bash
head -5 skills/pgedge-compactor/knowledge-base.md
```

Expected: first line is `# Skill Compactor Knowledge Base`, no
`---` fences.

- [ ] **Step 4: Commit**

```bash
git add skills/pgedge-compactor/knowledge-base.md
git commit -m "feat: add pgedge-compactor knowledge base

Compression taxonomy, safe/preserve lists, tabularization
guide with legend technique, diagnostic signals, and
anti-patterns. Adapted from gm-apprentice skill-compactor
with coding/pgedge examples throughout."
```

---

### Task 2: Create the SKILL.md

The main skill file with frontmatter, process, classification
taxonomy, compression techniques, rules, and reporting format.

**Files:**
- Create: `skills/pgedge-compactor/SKILL.md`

- [ ] **Step 1: Write SKILL.md**

Create `skills/pgedge-compactor/SKILL.md` with the following
content. The frontmatter `description` field is the trigger —
it must be a single long string covering all activation phrases
and exclusions.

```markdown
---
name: pgedge-compactor
description: "Use when the user wants to compress, compact, or optimize a Claude Code skill file (.md) to reduce token cost without losing quality-critical content. Trigger on phrases like 'compact this skill', 'compress the skill file', 'reduce token usage in this skill', 'optimize this SKILL.md', or when the user mentions a skill file being too long or too verbose. Also trigger when dispatched as a subagent with a target path for compaction. Do NOT trigger for general markdown editing, documentation writing (use pgedge-docs), or content that is not a Claude Code skill file."
---

# pgEdge Skill Compactor

You compact Claude Code skill files by compressing routing prose
while preserving priming content, concrete examples, and
cross-routing prompts.

## Before You Start

Read `knowledge-base.md` in this directory. It contains the
compression taxonomy and patterns learned from prior work.

## Target Resolution

If given a **single file**, use it directly.

If given a **directory**, use two-pass resolution:

1. Find `.md` files with YAML frontmatter containing `name:`
   and `description:` fields (skill entry points).
2. Scan those files for markdown links and "read X" references
   to other `.md` files in the same directory tree. Include
   those as targets too — they are part of the skill's token
   budget even though they lack frontmatter.
3. Exclude repo-level files: `README.md`, `CHANGELOG.md`,
   `LICENSE.md`.

If given a **glob**, expand and process each match.

If given **nothing**, ask the user for a target.

## Process

For each target file:

1. **Read the file in full.**

2. **Classify every block** as one of:

   | Category | Description | Action |
   |---|---|---|
   | Routing | Where to look: connective tissue, intros, "see X for Y", transitions | COMPRESS |
   | Priming | Shapes output: named patterns, thresholds, numbers, quality criteria, checklists | PRESERVE |
   | Examples | Concrete code samples, queries, usage demos, correct-output demonstrations | SACRED |
   | Tables | Already dense structured data | PRESERVE structure |
   | Cross-routing | Pointers to other files or skills | PRESERVE, strengthen |
   | Frontmatter | YAML between `---` fences | NEVER TOUCH |

3. **Compress routing blocks:**
   - Multi-sentence to single sentence.
   - Verbose bullets to terse bullets.
   - Redundant phrasing collapsed to single statement.
   - Introductory paragraphs to one line or removed.
   - Multi-paragraph descriptions to single-line directives
     with compact priming retained.
   - Parallel structured content (3+ items with consistent
     fields) converted to tables. See the Tabularization
     section in `knowledge-base.md` for guidance on when to
     tabularize and when not to, including the legend technique
     for repeated values.

4. **Verify preservation:**
   - All named patterns/frameworks still present with names.
   - All specific thresholds/numbers unchanged.
   - All concrete examples untouched.
   - All cross-routing prompts present.
   - All tables intact.
   - All trigger phrases verbatim.
   - Frontmatter byte-identical.

5. **Edit in place.** Surgical edits — one section at a time,
   not a full rewrite. Diffs must be reviewable.

6. **Report per file:**
   - Sections compressed (brief description of what changed).
   - Sections preserved (with reason: priming / example /
     table).
   - Before/after line count.
   - Judgment calls flagged for review.

After all files: report total files processed, total line
reduction, and any files that resisted compaction.

## Rules

- **Never compress to hit a target.** Compress what's safe.
  Report the result. A file that barely shrinks was already
  well-structured.

- **Never remove cross-routing prompts.** Pointers to other
  files/skills are the biggest quality lever. Strengthen if
  anything.

- **Never genericize domain-specific terminology.** If a skill
  uses specific terms (e.g., "Spock" not "replication
  extension", "SELECT_FIELD_SX" not "field styles"), preserve
  them verbatim.

- **Never touch frontmatter.** YAML between `---` fences is
  byte-identical before and after.

- **Deduplication needs breadcrumbs.** When content is
  consolidated, every location that previously had it gets a
  redirect: "See [topic] in [file]." Agents that find nothing
  stop looking.

- **Surgical edits only.** One section at a time, not full
  rewrites. Diffs must be reviewable.

- **Do not commit.** The user (or orchestrator) reviews the
  diff and commits if satisfied.

- **Report honestly.** If a file resists compaction (mostly
  tables, already terse), say so. Valid outcome, not failure.
```

- [ ] **Step 2: Verify frontmatter parses correctly**

Check that the frontmatter has valid `name` and `description`
fields and that no stray characters broke the YAML:

```bash
head -4 skills/pgedge-compactor/SKILL.md
```

Expected output:
```
---
name: pgedge-compactor
description: "Use when the user wants to compress, compact, or optimize a Claude Code skill file (.md) to reduce token cost without losing quality-critical content. Trigger on phrases like 'compact this skill', 'compress the skill file', 'reduce token usage in this skill', 'optimize this SKILL.md', or when the user mentions a skill file being too long or too verbose. Also trigger when dispatched as a subagent with a target path for compaction. Do NOT trigger for general markdown editing, documentation writing (use pgedge-docs), or content that is not a Claude Code skill file."
---
```

- [ ] **Step 3: Verify cross-routing to knowledge-base.md**

Confirm the SKILL.md references the knowledge base file:

```bash
grep -n "knowledge-base.md" skills/pgedge-compactor/SKILL.md
```

Expected: at least two matches — the "Before You Start" section
and the tabularization cross-routing in step 3.

- [ ] **Step 4: Commit**

```bash
git add skills/pgedge-compactor/SKILL.md
git commit -m "feat: add pgedge-compactor skill

Skill-file compactor that classifies content blocks as
routing/priming/examples/tables/cross-routing/frontmatter
and compresses only routing blocks. Supports single file,
directory (two-pass resolution), and glob targets."
```

---

### Task 3: Update README.md

Add the pgedge-compactor entry to the Skills section, following
the existing pattern (description bullets, invocation command,
example prompts).

**Files:**
- Modify: `README.md:59-98` (after pgedge-psql, before
  pgedge-webapp, or after pgedge-webapp before Support)

- [ ] **Step 1: Add pgedge-compactor section to README.md**

Insert the following after the `pgedge-webapp` section (after
line 99, before the `## Support & Resources` section). Follow
the same structure as existing entries:

```markdown
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
```

- [ ] **Step 2: Verify the README structure**

Check the README has the new section in the right position:

```bash
grep -n "^### " README.md
```

Expected: `pgedge-docs`, `pgedge-psql`, `pgedge-webapp`,
`pgedge-compactor` in that order, all before `## Support`.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: add pgedge-compactor to README"
```

---

### Task 4: Bump version in plugin.json and marketplace.json

Minor version bump (1.2.4 → 1.3.0) for new feature, no
breaking changes.

**Files:**
- Modify: `.claude-plugin/plugin.json:4` (version field)
- Modify: `.claude-plugin/marketplace.json:12` (version field)

- [ ] **Step 1: Update plugin.json version**

In `.claude-plugin/plugin.json`, change:

```json
"version": "1.2.4",
```

to:

```json
"version": "1.3.0",
```

- [ ] **Step 2: Update marketplace.json version**

In `.claude-plugin/marketplace.json`, change:

```json
"version": "1.2.4",
```

to:

```json
"version": "1.3.0",
```

- [ ] **Step 3: Verify versions match**

```bash
grep '"version"' .claude-plugin/plugin.json .claude-plugin/marketplace.json
```

Expected: both show `"version": "1.3.0"`.

- [ ] **Step 4: Commit**

```bash
git add .claude-plugin/plugin.json .claude-plugin/marketplace.json
git commit -m "chore: bump version to 1.3.0 for pgedge-compactor"
```

---

### Task 5: Final verification

Verify the complete skill is properly structured and all
cross-references resolve.

**Files:**
- Read: `skills/pgedge-compactor/SKILL.md`
- Read: `skills/pgedge-compactor/knowledge-base.md`
- Read: `README.md`

- [ ] **Step 1: Verify skill directory contents**

```bash
ls -la skills/pgedge-compactor/
```

Expected: `SKILL.md` and `knowledge-base.md`, nothing else.

- [ ] **Step 2: Verify SKILL.md references knowledge-base.md**

```bash
grep -c "knowledge-base.md" skills/pgedge-compactor/SKILL.md
```

Expected: 2 or more references.

- [ ] **Step 3: Verify no TTRPG/gm-apprentice references leaked in**

```bash
grep -ri "gm-apprentice\|ttrpg\|GM-guidance\|NPC\|ruling example\|FitD\|CoC\|GURPS\|BRP\|session-planner\|spotlight floor\|Five-Stage Arc\|Chekhov Protocol\|Three Clue Rule" skills/pgedge-compactor/
```

Expected: no matches.

- [ ] **Step 4: Verify README lists all four skills**

```bash
grep -c "^### pgedge-" README.md
```

Expected: 4.

- [ ] **Step 5: Verify versions are consistent**

```bash
grep '"version"' .claude-plugin/plugin.json .claude-plugin/marketplace.json
```

Expected: both show `1.3.0`.

- [ ] **Step 6: Review the full git log for this work**

```bash
git log --oneline -5
```

Expected: four commits — knowledge base, SKILL.md, README
update, version bump.
