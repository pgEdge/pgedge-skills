# Phase 1 — Intake and Verification

Nothing gets scored, drafted, or written until the vulnerability has been confirmed to exist, in real code, in real released versions. This phase produces the evidence every later phase cites.

## Step 1: Locate the Records Repository

This skill reads and writes a **records repository**: wherever this organization keeps its vulnerability disclosure process and its assessment records. There is no default path. If the user has not already named it in this conversation, ask before doing anything else — do not guess a location from a prior session, a similarly-named directory, or a convention from a different project.

Once established, every later reference in this skill to "the records repository" means this location. Do not write it into any file this phase or a later phase produces except the ones already living inside that repository — a drafted advisory, an assessment record, or anything else this skill writes for the user should not itself hardcode the path back to where it lives, for the same reason this skill does not.

## Step 2: Establish the Reporter

Determine, before anything else, whether this is an **internal** finding or an **external** report. Ask if it is not stated. The two paths diverge immediately:

| | Internal | External |
|---|---|---|
| Acknowledgement clock | none | starts at the moment the report arrived |
| Embargo | none by default | a default window applies unless a date is agreed |
| Credit | usually none; findings are pgEdge's own | the reporter is credited, in the form they asked for |
| Publication timing | tied to the release | tied to the agreed or default disclosure date |

For an external report, capture and record these before moving on:

- **The exact date and time the report was received.** Both clocks are computed from it. Do not estimate it, and do not infer it from the file timestamp of a forwarded email — ask the user for the received date and time if the report does not carry it unambiguously.
- **Who the reporter is**, as they would want to be named. A VDP team, a company's security organisation, and a named individual researcher are three different credits.
- **Whether a second party also deserves credit.** A fix commit may name an upstream researcher who found the same defect in a project pgEdge's code descends from. That person did not report to pgEdge, but the weakness is theirs. Ask the user how to handle it rather than deciding.
- **Any embargo terms the reporter proposed.** A reporter proposing a date the default policy would also produce is not a disagreement to resolve, but it still has to be confirmed back to them explicitly.

## Step 3: Read the Live Process Docs

Read these now, in full, from the records repository. Do not skip one because you believe you know what it says.

- **`process/DISCLOSURE-PROCESS.md`** — the disclosure process, the CVE determination tests, the timing tables, and the assessment-record template. This is the governing document.
- **`process/FILING-GUIDE.md`** — the mechanics of filing a GitHub Security Advisory: which file holds which value, the form's field order, and the rules that do not change.
- **A record already marked published, as a style template.** Identify one **by name**, from `records/README.md`'s outcome column only, then open **only that specific record's files** — `records/<that-one-name>/assessment-record.md`, `advisory.md`, `cve-record.json`. Do not grep, glob, or otherwise search across the content of `records/*` looking for a published one: doing so opens every record's files indiscriminately, including ones still under embargo, and a search that touches an unpublished record's content has already reached it, whether or not you go on to read the match. The README's outcome column is the only thing safe to search broadly, because it is metadata, not the confidential content. If the README's own outcome text does not make a record's published status unambiguous, ask the user rather than opening the record to check. If none of the available records happen to be published yet, work from `templates/` instead of an unpublished real one.
- **`records/README.md`** — the index of every record, its subject, and its outcome.
- **`governing-decisions.md`.** One table, one row per ruling on how the process applies to a recurring hard case — a partially masked credential fragment, whether publication discharges customer notification — held once so each new record does not re-argue settled ground. Read the whole table now; do not wait to see if a row turns out relevant, because the resemblance to a governing decision is often not obvious until partway through scoring.
- **`tools/form-fields.py`** — the CVE JSON schema as this program actually uses it, and the helper that renders a record into GitHub form fields.

If the user points at a handoff or status document as context, read it, but check its date. Those documents go stale; the procedure and filing guide govern where they disagree.

## Step 4: Locate and Pull the Product Repository

The vulnerability is in a product repository, which is separate from the records repository. Establish which one and where it is checked out locally — ask if it is not obvious from the report.

**Pull it from origin before reading anything.** A stale local checkout will show you code that has already been fixed, or hide a branch that has since been cut.

```bash
git -C /abs/path/to/product fetch --all --tags --prune
git -C /abs/path/to/product pull
```

Always address the repository with `git -C /abs/path`. A `cd` persists between tool calls and silently redirects every later command.

**Read only.** This skill never commits, pushes, branches, or tags in a product repository. If the user has supplied engineering's proposed fix as patch files, read them where they are; do not apply them to the tree.

## Step 5: Verify the Vulnerable Code Exists

Do not take the report's line citations on faith. Reports cite the version the reporter looked at, which may not be the version that shipped, and line numbers drift between releases.

1. **Find the function or code path the report names**, by searching the actual source rather than navigating to the cited line number.
2. **Read the surrounding code and confirm the described pattern is genuinely there** — the specific unchecked value, the specific missing validation, the specific call. Write down what the pattern actually is, in your own words, at file and line.
3. **Confirm it is unchanged**, not already fixed by an unrelated commit since the reporter looked.

Record the result as evidence with file and line, per tag, exactly as the assessment record expects. "Confirmed by direct inspection on *date*" with the concrete pattern quoted is the standard; "the report says so" is not.

Where the report claims a full exploit chain that pgEdge has not independently re-run — a ROP chain to code execution, for instance — say so explicitly, and separate the part that *is* independently confirmed (the underlying memory-safety or logic defect) from the part taken on the reporter's demonstration. Both statements belong in the record.

## Step 6: The Multi-Branch Trap

**Do not assume the product has one "latest release."** Some pgEdge products support several major versions concurrently — a support policy of N years from each major's first release puts three or four majors in the window at once. Others really do have a single line. Which one you are looking at changes the affected-versions section, the CVE JSON structure, and the number of patched versions engineering has to commit to.

Establish the product's actual support policy first — from the product's own documentation or by asking the user. Then verify the vulnerable pattern in the **latest GA tag of each supported line**, not at HEAD alone.

```bash
# What lines exist, in release order
git -C /abs/path/to/product tag --sort=creatordate

# The latest tag on each major line
git -C /abs/path/to/product tag --list 'v5.*' --sort=-v:refname | head -5

# The vulnerable file as it stands in a given tag
git -C /abs/path/to/product show v5.0.11:src/path/to/file.c | less

# Or search the file at that tag directly
git -C /abs/path/to/product show v5.0.11:src/path/to/file.c | grep -n 'function_name'
```

For each supported line, record: the earliest affected tag, the latest GA tag, and the file and line where the pattern appears **in that tag**. Line numbers differ between tags for the same defect; record each one rather than one number for all of them.

Include pre-release lines that are downloadable. A major still in beta is affected if the code is there, and its users can get it.

This is also the section where the CVE released-version test is really decided. A defect present only at HEAD and never in a released tag reaches a different outcome — see Phase 3.

**A tag with no branch ancestry to anything current is not evidence that its line is retired.** Some products carry an old release line whose branch was deleted or renamed after release, so the tag survives as an orphan with no path back to `main` or any live branch — do not read that absence of ancestry as absence of support. Check directly:

```bash
# Does any current branch contain this tag's commit?
git -C /abs/path/to/product branch --all --contains v4.1.0.3

# If none do, the tag is orphaned - check the product's own published
# support lifecycle (a docs site's version navigation, a release
# policy page) for whether that line is still supported, rather than
# inferring retirement from git topology alone
```

Where the product's own documentation does not resolve it either, say so as an open question rather than silently assuming retired — it changes whether that line needs a patched version at all, and guessing wrong in either direction is a defect the record should surface, not absorb.

## Step 7: Scope the Report

A report describes one weakness. An engineering fix may address several.

If the fix you have been given bundles more than one distinct issue — a code-execution defect plus some denial-of-service hardening plus a memory-leak cleanup, say — **scope this advisory to what was actually reported and scored, and flag the rest to the user as separate future work.** Do not silently bundle unrelated issues of different severity into one advisory, and do not silently drop them either. The procedure document's rule that one report describing several weaknesses gets a record each applies here in reverse: one fix addressing several weaknesses does not get one advisory.

State the scoping decision, and what it excludes, in writing. Phase 3 records it.

## Output of This Phase

Before moving to Phase 2, you should be able to state:

- the records repository's location
- internal or external; if external, reporter identity, received date and time, and credit preference
- the product, the repository, and which of its lines are currently supported
- the vulnerable pattern, in your own words, with file and line per supported tag, verified by inspection
- what is independently confirmed versus what rests on the reporter's demonstration
- what this advisory covers and what it deliberately excludes

Proceed to Phase 2.
