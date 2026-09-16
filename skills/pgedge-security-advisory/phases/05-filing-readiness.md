# Phase 5 — Filing Readiness

This skill produces drafts. It does not file them.

## The Hard Limits

These are not defaults to be overridden by a user saying "go ahead". They are the shape of the skill.

- **Never request a CVE identifier.** Identifiers are never reused, and withdrawing one is a public rejection. An identifier requested for an advisory that is then not filed sits in RESERVED permanently.
- **Never press Publish, and never instruct anyone to.** Requesting an identifier discloses nothing; publishing discloses everything, and it is irreversible.
- **Never create the draft advisory on GitHub.** Even a draft advisory is a state change in a product repository's security tab.
- **Never touch the product repository.** No commits, no branches, no pushes, no tags — and in particular, never push a fix that is being deliberately held out of the public repository while under embargo. A security fix pushed to a public branch under a descriptive commit message *is* the disclosure, whether or not anyone intended it.

If the user asks this skill to do one of these things, say plainly that it will not, and hand back the checklist instead. Filing is a human action taken with the checklist in hand.

## The Handback Checklist

Produce a numbered list of everything still open, specific to this advisory, with the current state of each item rather than a generic template. Every item is either **clear** (with evidence and a date) or **open** (with who has to close it).

Work through at least these, adding anything Phase 3 flagged:

1. **Has the acknowledgement gone out?** External reports only. State the deadline, the arithmetic behind it, and whether it has passed. If it has not been sent, this is the first item and nothing else matters yet.
2. **Has the disclosure date been confirmed back to the reporter?** Not what the policy default produces, and not what the reporter proposed — what has actually been agreed in writing, by whom, on what date.
3. **Is there a patched version, and is it real?** Per supported branch. Distinguish confirmed from provisional, and name who confirms each provisional one. A version named ahead of its release is acceptable in a draft and unacceptable at publication.
4. **Is the fix actually merged, and where?** A fix held as unmerged patch files under embargo is a different state from a fix on `main`, and it changes both the disclosure position and what the advisory can claim.
5. **Who gates publication for this product, and have they signed off?** Never assume the approval chain established for other products extends to this one. If this has not been confirmed, it is an open item, not a formality.
6. **Does any customer hold an individual notification obligation for this product?** Assessed at publication, not when the record was written.
7. **Has the base score been re-confirmed against GitHub's own calculator?** Hand-computed at Phase 2; verified when the vector is pasted into the form.
8. **Has the assessment record been signed off?** A record with no name and date against sign-off does not serve as audit evidence.

Close with the button order from `FILING-GUIDE.md` — fill the fields, create the draft, request the identifier on the saved draft, stop — so that whoever files does not have to go looking for it. And state explicitly which of those steps is the one nobody takes without the approvals above.

## Committing the Drafts

**Only when explicitly asked.** Drafting the files is the deliverable; committing them is a separate request.

When asked, commit in the records repository only, at the location established in Phase 1:

```bash
git -C <records-repo> add \
  records/AR-YYYY-NNN-slug/assessment-record.md \
  records/AR-YYYY-NNN-slug/advisory.md \
  records/AR-YYYY-NNN-slug/cve-record.json \
  records/README.md

git -C <records-repo> commit
```

- **Stage explicit paths. Never `git add -A` or `git add .`.** A broad stage sweeps unrelated in-progress edits into a security commit.
- **Always `git -C /abs/path`.** A `cd` earlier in the session silently redirects a bare `git commit` into whichever repository was last visited, which here could be the product repository.
- Conventional commit style, and no self-attribution.
- The commit message says what was drafted and that it is not filed. `docs: draft <product> <weakness> advisory (AR-YYYY-NNN)` is the established shape.
- **Commit only. Never push, without being separately asked to push.** Committing locally and pushing to the shared remote are different requests — this repository has a remote other people read, unlike the working copy it replaced.

Committing a draft does not change the filing position. The advisory remains unfiled and unrequested until a human takes the steps in `process/FILING-GUIDE.md`.
