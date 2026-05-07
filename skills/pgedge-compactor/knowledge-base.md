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
