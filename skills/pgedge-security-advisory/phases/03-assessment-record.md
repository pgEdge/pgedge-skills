# Phase 3 — Assessment Record

The assessment record is the audit trail. It records why every decision went the way it did, with evidence, and it is what a reviewer or an auditor reads instead of asking you.

It lives at `assessment-record.md` inside the record's own folder under the records repository's `records/` directory.

## Step 1: Read the Template Fresh

The record's field list is defined in `process/DISCLOSURE-PROCESS.md`, in the section describing the record itself. **Read it now and transcribe the field list from what is actually there.** Do not reconstruct it from memory, from this file, or from an older record — the template has changed before and the procedure document governs.

`templates/assessment-record-template.md` in this skill is a structural scaffold only. Use it for the shape of the file; take the field list from the live procedure document, and reconcile the two if they differ, in the live document's favour.

Then read the two or three most recent records under `records/` for the house style: a fenced block carrying the structured fields, followed by prose sections for anything the fields cannot hold.

## Step 2: Number It and Create the Folder

Records are numbered `AR-YYYY-NNN`. Take the next unused number:

```bash
ls <records-repo>/records/
```

Substitute the records repository's actual location, established in Phase 1, for `<records-repo>` in this and every other command in this skill.

**Never reuse a number**, including one belonging to a record that closed no-CVE. Gaps are acceptable; collisions are not. Create `records/AR-YYYY-NNN-short-slug/`, and write the assessment record to `assessment-record.md` inside it. The same slug names the folder, the CVE JSON record, and everything else about this vulnerability — there is no separate numbering scheme to reconcile it against.

## Step 3: Fill the Evidence Sections

Everything Phase 1 verified goes here, with dates.

- **Reachable as shipped** — what was confirmed, by whom, and how. Separate what pgEdge independently verified from what rests on the reporter's demonstration. If a regression test pins the behaviour, name the test file.
- **Affected versions** — per supported line: the earliest affected tag, the latest GA tag, and the file and line where the pattern appears **in that tag**. The procedure document is explicit that this is where assessments go wrong: test against the release tags, not against the fix's commit date.
- **Fixed in** — the patched version, or `UNRELEASED` with the state of the fix. Where a version number is named ahead of its release, mark it clearly as provisional and name who has to confirm it. Where one line's version is confirmed and others are not, say which is which; do not let a single confirmed number imply the rest.
- **Severity** — the vector, the score, and the judgement-call reasoning from Phase 2, including the case not taken.

## Step 4: Compute the Disclosure Clock

For an external report only. Take the limits from the timing table in the live procedure document — **read the numbers there, do not carry them in your head**, they have been changed before.

Two points that are easy to get wrong:

- **There are two distinct clocks and they do not start at the same event.** One runs from identification and governs how long remediation may take. The other runs from fix availability and governs communication. A vulnerability can be comfortably inside its remediation window with nothing yet to communicate. Read the procedure document's own description of which is which, and do not conflate them, or quote one as if it were the other.
- **The acknowledgement deadline is in business days.** Compute it by listing the dates explicitly — write the days out in the record so the arithmetic is auditable — and flag it loudly if it is imminent or already passed. An acknowledgement that has not gone out is the most urgent thing in the file.

Record the publication date as agreed, or as the default the policy produces, and note explicitly whether it has been **confirmed back to the reporter** or is merely what the policy would give. Those are different states and only one of them is an agreement.

Note the precedence order the procedure document sets out for conflicting obligations, and apply it rather than defaulting to the table.

## Step 5: Check Governing Decisions Before Running the Tests

Before applying the CVE-eligibility test or deciding the disclosure and notification path, check whether a ruling in `governing-decisions.md` already covers this shape of finding — Phase 1 read all of them for exactly this reason. Where one applies, cite it by its `AR-YYYY-NNN` number and follow it rather than re-deriving the same conclusion. Where a governing decision names a limit on when it applies (most of them do), confirm this record actually falls inside that limit before relying on it — do not cite a governing decision from memory of its headline ruling alone.

If this record's shape resembles a governing decision closely enough that a future record would benefit from the same ruling being held once, say so as a recommendation in the record rather than writing a new governing decision yourself. Creating one is a deliberate step the user takes, not something this phase does on its own.

## Step 6: Run the CVE-Eligibility Tests

The procedure document defines the tests and the outcome table. **Read them there.** As of writing there are two, both of which must pass, and the document carries a table of situations that fail one or the other — but confirm the current wording rather than assuming this description is still accurate.

Record each test as pass or fail **with its evidence**, not as a bare verdict:

- For the released-version test, the evidence is the tag list from Phase 1 — the vulnerable code confirmed present in released GA tags, not merely at HEAD or in a pre-release.
- For the trust-boundary test, the evidence is a concrete identification of the untrusted party and what privilege they hold. Name them: "a database role holding CREATE on a replicated schema" is evidence; "an attacker" is not.

Then record the outcome and, for the identifier field, `NOT YET REQUESTED` with the specific preconditions that must clear first. Never request an identifier in this phase or any other.

## Step 7: Flag Unresolved Process Questions — Do Not Resolve Them

Where the existing process does not obviously cover this case, **write the gap into the record as an open question** and carry it forward to Phase 5's checklist. Do not work around it, and do not assume it resolves in the convenient direction.

The gaps worth checking for on any new advisory:

- **Does the existing publish-approval chain cover this product?** The approvals this program uses were established for a particular set of products. A vulnerability in a different product may or may not fall under the same approvers. Confirming who gates publication for *this* product is a checklist item, never an assumption.
- **Does any customer hold an agreement carrying an individual notification obligation for this product?** `AR-2026-018` in `governing-decisions.md` rules that publication is the notification channel for everyone — read its stated limit before treating this as settled for a customer whose contract has not been checked against that limit. A prior record's finding of "no such customers" was scoped to that product's contracts and does not transfer.
- **Does existing tooling handle this record's shape?** See Phase 4 on the multi-branch case.
- **Is there a precedent record for this situation at all?** The first externally-reported issue, the first multi-branch issue, and the first issue in a new product each break new ground. Say so at the top of the record, so the next reader knows there was no pattern to match against.

## Step 8: Add Prose Sections for What the Fields Cannot Hold

The fenced field block is fixed; anything else goes below it as a prose section with its own heading. Common ones, drawn from existing records:

- **Scoping** — what this record covers and what it deliberately excludes, when the fix addresses more than was reported.
- **Customer action beyond upgrading** — log review, credential rotation, treating old output as sensitive. Upgrading stops further exploitation; it does not tell anyone whether the window was already used, and readers will not infer the step.
- **Workaround** — and honestly, what it does not cover.
- **Known tooling or process gaps** affecting this record specifically.

## Step 9: Update the Index

`records/README.md` carries a table indexing every record. Add a row: record number, one-line subject, and outcome (`CVE, 7.5 High`, `No CVE, closed on evidence`, and so on). Mark a record that is drafted but not filed as such.

If the README has narrative sections that this record contradicts or extends — a statement about which reports were internal, say, or about who the customer-notification class covers — flag the tension to the user rather than editing the narrative silently.

## Output of This Phase

A complete `records/AR-YYYY-NNN-slug/assessment-record.md`, an updated index row, and a list of open process questions carried forward.

Proceed to Phase 4.
