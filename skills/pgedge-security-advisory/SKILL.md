---
name: pgedge-security-advisory
description: "Use when a security vulnerability is reported for a pgEdge product — an internal finding, a vulnerability disclosure email, or a report from an external researcher or VDP program — or when asked to investigate, score (CVSS), or draft/file a security advisory, assessment record, or CVE record for a pgEdge repository. Also use when asked to check which supported versions or branches of a product are affected by a vulnerability. Do not trigger for general security code review, dependency-scan triage, or fixing an already-ticketed bug with no disclosure component."
---

# pgEdge Security Advisory

Take a reported vulnerability from raw report to a drafted assessment record, advisory draft, and CVE JSON record — investigated, scored, and ready for a human to file.

## The Records Repository — Location Is Never Hardcoded

This skill reads from and writes to a **records repository**: wherever your organization keeps its vulnerability disclosure process and its assessment records. This skill carries no fixed path to one, on purpose — a skill shipped publicly must not bake in a pointer to anyone's private infrastructure, and different installs keep this repository in different places. Phase 1 establishes its location at the start of every run; ask the user for it if it is not already clear from context. Every later phase refers to "the records repository", never a literal path.

## The Live Policy Docs Are Authoritative — Read Them Every Run

`process/` in the records repository is the source of truth for every policy specific: the acknowledgement and embargo clocks, the CVE determination tests, the assessment-record template, and the filing mechanics.

**Read those documents fresh at the start of every run. Never work from this skill's memory of them, and never from your own.** They change: patched-version numbers have been corrected, the acknowledgement window has been rewritten, the CVE determination test has been reworded. Anything this skill says about *policy* is a description of where to look, not a copy of what you will find.

What this skill owns permanently is the process discipline that does not change — and `standards/CVSS-JUDGMENT-CALLS.md`, which is scoring judgement, not policy.

## Core Principles

These hold on every run, regardless of what the policy docs say.

1. **Verify independently. Never trust the report.** A report's line citations, version claims, and severity assessment are inputs to check, not facts to transcribe. Read the real source in the real repository.
2. **Check every supported branch, not just the latest.** Some pgEdge products support several major lines concurrently. Assuming "latest release only" produces an advisory that tells most of the affected population they are safe.
3. **Ask, do not silently decide.** Where a genuine judgement call exists — a CVSS metric that reasonable people score either way, the scope of an advisory against a multi-part fix, how a reporter wants to be credited — put the question to the user in plain language with the consequence of each answer. Do not pick one and move on.
4. **Never file, never request a CVE identifier, never publish.** This skill produces drafts and a readiness checklist. A human presses every button.

## Phase Sequence

### Phase 1 — Intake and Verification

Read `phases/01-intake-and-verification.md` and follow its instructions.

Establish who reported it and whether they are internal or external. Read the live process docs and every governing decision. Locate and pull the affected product repository. Confirm the vulnerable code exists, in the latest tag of every currently-supported line.

Do not score anything during this phase.

### Phase 2 — CVSS Scoring

Read `phases/02-cvss-scoring.md` and follow its instructions.

Build the vector one metric at a time with stated reasoning. Escalate genuinely ambiguous metrics to the user rather than choosing. Hand-verify the base score. Choose CWEs.

`standards/CVSS-JUDGMENT-CALLS.md` carries the worked reasoning for the two metrics that most often need a judgement call.

### Phase 3 — Assessment Record

Read `phases/03-assessment-record.md` and follow its instructions.

Write the record into a new folder under the records repository's `records/` directory, following the template as the live procedure document currently defines it. Check `governing-decisions.md` for an applicable precedent before running the CVE-eligibility tests. Compute the disclosure clock, and flag unresolved process questions rather than resolving them silently.

### Phase 4 — Advisory Draft and CVE Record

Read `phases/04-advisory-and-cve-record.md` and follow its instructions.

Write the CVE JSON record as the schema source of truth, and the advisory draft as a field-for-field transcription of GitHub's advisory form, pasteable top to bottom.

### Phase 5 — Filing Readiness

Read `phases/05-filing-readiness.md` and follow its instructions.

Hand back a numbered checklist of what is still open before anyone files. Commit the drafts only if asked, only with explicit paths, only in the records repository, and never push without being separately asked.

## Templates

`templates/` holds structural skeletons for the three output files — placeholder tokens only. They save reconstructing file shape from scratch. They deliberately carry no policy text: field meanings live in the phase files, policy specifics live in the external docs.
