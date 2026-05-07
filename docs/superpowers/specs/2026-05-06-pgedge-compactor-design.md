# pgedge-compactor Skill Design

## Summary

A Claude Code skill that compacts skill files by compressing
routing prose while preserving priming content, concrete examples,
and cross-routing prompts. Reduces token cost without losing
quality-critical content.

Based on compression patterns learned from the gm-apprentice
skill-compactor project, generalized for any Claude Code skill
file with pgedge-specific examples throughout.

## File Structure

```
skills/pgedge-compactor/
  SKILL.md            # Process, classification, rules, invocation
  knowledge-base.md   # Compression taxonomy, patterns, anti-patterns
```

## Skill Metadata

- **Name:** `pgedge-compactor`
- **Trigger:** When the user wants to compress, compact, or
  optimize a Claude Code skill file to reduce token cost. Phrases
  like "compact this skill", "compress the skill file", "reduce
  token usage in this skill", "optimize this SKILL.md", or when
  a skill file is described as too long or too verbose. Does NOT
  trigger for general markdown editing, documentation writing
  (use pgedge-docs), or non-skill-file content.

## Invocation Model

- **Direct:** User invokes `/pgedge-skills:pgedge-compactor` and
  points it at a file, directory, or glob.
- **Subagent:** Other skills or orchestrators dispatch the
  compactor as a subagent with a target path.

## Target Resolution

- **Single file:** Use it directly.
- **Directory:** Two-pass resolution. First, find `.md` files
  with YAML frontmatter containing `name:` and `description:`
  fields (the skill entry points). Second, scan those files for
  markdown links and "read X" references to other `.md` files
  in the same directory tree and include those as targets too.
  This catches knowledge bases, reference files, and other
  routing targets that are part of the skill's token budget
  but lack frontmatter. Excludes repo-level files like
  `README.md`, `CHANGELOG.md`, `LICENSE.md`.
- **Glob:** Expand and process each match.
- **Nothing provided:** Ask the user.

## Core Process

### Step 0 - Load Knowledge Base

Read `knowledge-base.md` from the skill directory on activation.

### Step 1 - Identify Targets

Resolve the input to a list of files using the target resolution
rules above.

### Step 2 - Per-File Compaction Loop

For each target file:

1. **Read the file in full.**

2. **Classify every block** using the taxonomy:

   - **Routing** - tells the agent where to look. Connective
     tissue, introductory prose, "see X for Y", section
     transitions. COMPRESSIBLE.
   - **Priming** - shapes output quality. Named patterns,
     thresholds, specific numbers, quality criteria, checklist
     items. PRESERVE.
   - **Examples** - concrete code samples, query examples, usage
     demonstrations, anything that shows what correct output
     looks like. SACRED - do not touch.
   - **Tables** - already dense. Compress prose around them, not
     the tables themselves. PRESERVE structure.
   - **Cross-routing** - pointers to other files or skills.
     PRESERVE and potentially strengthen.
   - **Frontmatter** - YAML between `---` fences. NEVER TOUCH.

3. **Compress routing blocks only:**

   - Multi-sentence to single sentence.
   - Verbose bullets to terse bullets.
   - Redundant phrasing collapsed to single statement.
   - Introductory paragraphs reduced to one line or removed.
   - Multi-paragraph descriptions to single-line directives
     with compact priming retained.
   - Parallel structured content (repeated bullet patterns,
     if/then lists, option descriptions with consistent fields)
     converted to tables when the content has 3+ items sharing
     the same structure. Tables preserve all information at
     significantly lower token cost.

4. **Verify preservation:**

   - All named patterns/frameworks still present with names.
   - All specific thresholds/numbers unchanged.
   - All concrete examples untouched.
   - All cross-routing prompts present.
   - All tables intact.
   - All trigger phrases verbatim.
   - Frontmatter byte-identical.

5. **Edit in place.** Surgical edits - one section at a time,
   not a full rewrite. Diffs must be reviewable.

6. **Report per file:**

   - Sections compressed (brief description of what changed).
   - Sections preserved (with reason: priming / example / table).
   - Before/after line count.
   - Judgment calls flagged for review.

### Step 3 - Summary

After all files: total files processed, total line reduction,
any files that resisted compaction (already well-structured).

## Rules

- **Never compress to hit a target.** Compress what's safe.
  Report the result. A file that barely shrinks was already
  well-structured.
- **Never remove cross-routing prompts.** Pointers to other
  files/skills are the biggest quality lever. Strengthen if
  anything.
- **Never genericize domain-specific terminology.** If a skill
  uses specific terms (e.g., "Spock" not "replication extension",
  "SELECT_FIELD_SX" not "field styles"), preserve them verbatim.
- **Never touch frontmatter.** YAML between `---` fences is
  byte-identical before and after.
- **Deduplication needs breadcrumbs.** When content is
  consolidated, every location that previously had it gets a
  redirect: "See [topic] in [file]." Agents that find nothing
  stop looking.
- **Surgical edits only.** One section at a time, not full
  rewrites. Diffs must be reviewable.
- **Do not commit.** The user (or orchestrator) reviews the diff
  and commits if satisfied.
- **Report honestly.** If a file resists compaction (mostly
  tables, already terse), say so. Valid outcome, not failure.

## Knowledge Base Structure

The `knowledge-base.md` file contains:

### The Taxonomy

The classification system with coding/pgedge examples:

1. **Routing compresses** - text that tells the agent where to
   look. "Read knowledge-base.md in this directory" is routing.
2. **Priming preserves** - text that shapes output quality.
   "Enforce 90% coverage gate. Use `--csv -t` for queries." is
   priming.
3. **Concrete examples are sacred** - code samples, query
   examples, and usage demonstrations teach the agent what
   correct output looks like. Removing the psql format
   comparison table and flag examples caused the agent to default
   to verbose aligned output.
4. **Deduplication needs breadcrumbs** - when content lives in
   one canonical location, every file that used to have it needs
   a redirect, not silence.

### What Compresses Safely

- Introductory paragraphs ("This skill guides you through...").
- "See X for Y" prose (shorten, don't remove the pointer).
- Section transitions and connective tissue.
- Multi-sentence routing that can become single-line arrows.
- Redundant phrasing saying the same thing twice.
- Verbose bullet explanations where the bullet text is clear.
- Multi-paragraph command descriptions reduced to single-line
  directives with compact priming.

### What Must Be Preserved

- **Named patterns** - Form Field Pattern, Read-Only Access
  Pattern, Token-Efficient Output Formats, Keep a Changelog
  format, Safe Write Pattern.
- **Specific thresholds** - 90% coverage gate, 79-character line
  wrapping, port range 1024-65535, LIMIT 5 for initial
  exploration, 12-character minimum password.
- **Tables** - already near-optimal density. Compress prose
  around them, not the tables themselves.
- **Checklist items and coverage criteria** - quality gates the
  agent checks against.
- **Trigger phrases** - verbatim, these are routing anchors.
- **Frontmatter** - never touch YAML frontmatter.
- **Concrete examples** - code samples, SQL queries, config
  snippets, any content that demonstrates what correct output
  looks like for a specific domain.
- **Cross-routing prompts** - "For X, also read Y.md" - the
  biggest quality lever. Compacted files feel "complete enough"
  so agents stop; cross-routing nudges them to keep going.

### Tabularization

Converting parallel prose into tables is one of the highest-yield
compression techniques. A well-formed table preserves all
information at a fraction of the token cost. But it's easy to
misapply — the knowledge base needs clear guidance on when to
tabularize and when not to.

**Candidates for tabularization (3+ items required):**

- Lists of options, flags, or settings where each item has the
  same attributes (name, default, description, type).
- Comparison lists ("X does this, Y does that, Z does the
  other") where items share a consistent structure.
- Conditional behavior ("when A, do B; when C, do D") with
  consistent condition/action pairs.
- Format or mode descriptions with consistent fields. Example:
  the psql skill's format summary table was already tabular —
  Scenario | Flags | Why — and is a model of what good
  tabularization looks like.
- Error code or status code mappings.
- Permission or role matrices.

**Do NOT tabularize:**

- Sequential procedures (numbered steps depend on order and
  context from previous steps — tables lose that).
- Code blocks and concrete examples (sacred content — tables
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
- Cell content should be terse — if a cell needs more than ~10
  words, the content may not be a good table candidate.
- Preserve all information from the original prose. If
  tabularizing forces you to drop nuance, keep the prose.
- Add a one-line introductory sentence before the table
  (consistent with pgedge-docs style).

**Legend technique for repeated values:**

When the same long values appear across many rows, replace them
with short codes and add a legend above the table. This
compresses the table dramatically without losing information.

Example — a permissions matrix before compression:

| Endpoint | superuser | administrator | editor | viewer |
|---|---|---|---|---|
| GET /users | read-write | read-write | read-only | read-only |
| POST /users | read-write | read-write | none | none |
| DELETE /users | read-write | none | none | none |
| GET /clusters | read-write | read-write | read-only | read-only |

After applying the legend technique:

**su**: superuser, **adm**: administrator, **ed**: editor,
**vi**: viewer — **RW**: read-write, **RO**: read-only,
**--**: none

| Endpoint | su | adm | ed | vi |
|---|---|---|---|---|
| GET /users | RW | RW | RO | RO |
| POST /users | RW | RW | -- | -- |
| DELETE /users | RW | -- | -- | -- |
| GET /clusters | RW | RW | RO | RO |

The legend costs a few tokens once; the savings compound across
every row. Apply when 3+ values repeat across 4+ rows.

### Diagnostic Signals

- **File resists compaction** - already well-structured. The
  psql skill's format summary table and schema discovery SQL
  blocks are 80% structured content. Minimal savings is success.
- **Verbose prose dilutes signal** - after compaction, actionable
  content is a larger proportion of what the agent reads.
- **Tables compress beautifully** - they preserve information
  density at lower token cost than equivalent prose.

### Anti-Patterns

- Cutting priming to hit a size target (the routing/priming
  boundary is sacred).
- Removing concrete examples because they seem like "filler"
  (output quality drops when the agent loses calibration
  examples).
- Deduplicating without leaving redirects (agents find nothing
  and stop looking).
- Compressing domain-specific terminology into generic terms
  (replacing "psql" with "database CLI" caused the agent to
  suggest alternative tools - domain terms are routing anchors).
- Removing cross-routing prompts as redundant (they're the
  single biggest quality lever).

## Integration with pgedge-skills Repo

### Plugin Registration

No changes needed to `plugin.json` or `marketplace.json` - the
skill is auto-discovered from the `skills/` directory by the
Claude Code plugin loader.

### README Update

Add a `pgedge-compactor` section to `README.md` following the
same pattern as the existing skill entries.

### Version Bump

Bump the plugin version in `plugin.json` and `marketplace.json`
per semver (minor bump - new feature, no breaking changes).
