---
name: pgedge-compactor
description: "Use when the user wants to compress, compact, or optimize a Claude Code skill file (.md) to reduce token cost without losing quality-critical content. Trigger on phrases like 'compact this skill', 'compress the skill file', 'reduce token usage in this skill', 'optimize this SKILL.md', 'split this reference file', 'this reference file is too big to load every time', or when the user mentions a skill file being too long or too verbose. Also trigger when dispatched as a subagent with a target path for compaction. Do NOT trigger for general markdown editing, documentation writing (use pgedge-docs), or content that is not a Claude Code skill file."
---

# pgEdge Skill Compactor

You compact Claude Code skill files two ways: by compressing routing prose, and - when prose is exhausted or a file mixes always-needed with conditional content - by structurally extracting content into files that load only when needed. Both axes preserve priming content, concrete examples, and cross-routing prompts.

## Before You Start

Read `knowledge-base.md` in this directory. It contains the compression taxonomy and patterns learned from prior work.

## Target Resolution

If given a **single file**, use it directly.

If given a **directory**, use two-pass resolution:

1. Find `.md` files with YAML frontmatter containing `name:` and `description:` fields (skill entry points).
2. Scan those files for markdown links and "read X" references to other `.md` files in the same directory tree. Include those as targets too — they are part of the skill's token budget even though they lack frontmatter.
3. Exclude repo-level files: `README.md`, `CHANGELOG.md`, `LICENSE.md`.

If given a **glob**, expand and process each match.

If given **nothing**, ask the user for a target.

## Process

Compaction has two axes. Start with **Axis 1 — prose compression**. When a file's prose pass yields <1% (it is exhausted) or a file mixes always-needed content with content only some invocations use, move to **Axis 2 — structural extraction**. See `knowledge-base.md` for the full method behind both.

### Axis 1 — Prose compression

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
   - Multi-paragraph descriptions to single-line directives with compact priming retained.
   - Parallel structured content (3+ items with consistent fields) converted to tables. See the Tabularization section in `knowledge-base.md` for guidance on when to tabularize and when not to, including the legend technique for repeated values.

4. **Verify preservation:**
   - All named patterns/frameworks still present with names.
   - All specific thresholds/numbers unchanged.
   - All concrete examples untouched.
   - All cross-routing prompts present.
   - All tables intact.
   - All trigger phrases verbatim.
   - Frontmatter byte-identical.

5. **Edit in place.** Surgical edits — one section at a time, not a full rewrite. Diffs must be reviewable.

6. **Report per file:**
   - Sections compressed (brief description of what changed).
   - Sections preserved (with reason: priming / example / table).
   - Before/after line count.
   - Judgment calls flagged for review.

After all files: report total files processed, total line reduction, and any files that resisted compaction.

### Axis 2 — Structural extraction

When a file's prose pass is exhausted (<1%) or the file mixes always-needed content with content only some invocations use, cut the fixed per-invocation load by splitting the file. See the Structural Extraction section in `knowledge-base.md` for the full method, including when NOT to split.

1. **Confirm it applies.** Prose pass yielded <1%, or the file bundles independent topics / a conditional block. If not, stop - never split a small or always-loaded file.

2. **Split on natural seams.** Either *split-and-index* (one file per topic, thin index at the original path) or *pull-conditional* (move a self-contained block to its own file).

3. **Byte-preserve content.** Extraction moves bytes between files, never rewords them. Prose compression is a separate pass.

4. **Breadcrumb and re-point.** Leave a redirect at every old location; update the links, "read X" references, and tests that pointed at the moved content.

5. **Report** the per-invocation load before/after and the new file layout.

## Rules

- **Never compress to hit a target.** Compress what's safe. Report the result. A file that barely shrinks was already well-structured.

- **Never remove cross-routing prompts.** Pointers to other files/skills are the biggest quality lever. Strengthen if anything.

- **Never genericize domain-specific terminology.** If a skill uses specific terms (e.g., "Spock" not "replication extension", "SELECT_FIELD_SX" not "field styles"), preserve them verbatim.

- **Never touch frontmatter.** YAML between `---` fences is byte-identical before and after.

- **Deduplication needs breadcrumbs.** When content is consolidated, every location that previously had it gets a redirect: "See [topic] in [file]." Agents that find nothing stop looking.

- **Surgical edits only (Axis 1).** Prose compression works one section at a time, not full rewrites. Diffs must be reviewable.

- **Extraction preserves bytes (Axis 2).** Splitting a file moves content verbatim between files and leaves breadcrumbs; never reword during a split. Prose compression and extraction are separate passes, never the same edit.

- **Do not commit.** The user (or orchestrator) reviews the diff and commits if satisfied.

- **Report honestly.** If a file resists compaction (mostly tables, already terse), say so. Valid outcome, not failure.
