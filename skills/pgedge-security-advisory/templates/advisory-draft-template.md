# {SHORT TITLE}

Record: AR-{YYYY}-{NNN}. CVE record: cve-record.json, same folder.
Repository to draft in: pgEdge/{repo}

Draft at:
https://github.com/pgEdge/{repo}/security/advisories/new

Laid out {YYYY-MM-DD} to mirror the GitHub form: same field order, same field names, description already written to GitHub's own template (Impact / Patches / Workarounds / References). Fill it top to bottom without scrolling back and forth in this file.

{**DO NOT FILE THIS YET.** — with one line saying why. Delete both if nothing is outstanding, and see "Before you touch the form" below.}

---

## Field 1 — Title \*

{One unwrapped line. Name the product, the component, and the outcome.}

## Field 2 — CVE identifier

Select: **Request CVE ID later**

## Field 3 — Description \*

Paste everything between the two markers. Nothing above or below them.

<!-- vvvvvvvvvvvvvv  START COPYING ON THE NEXT LINE  vvvvvvvvvvvvvv -->

### Impact

{What an attacker can do, and what privilege they need to start. Lead with the privilege, because it is what lets a reader rule themselves in or out.}

{How the attack works, in enough detail to be credible and not enough to be a recipe.}

**Who is affected:** {which deployments, which configuration, which versions.}

{Any upstream or sibling-project relationship — the same defect fixed elsewhere, with its identifier.}

### Patches

{**Fixed in `{version}`. Upgrade to `{version}` or later.** — or **Not yet released.** with the state of the fix.}

{What the fix does, one short paragraph.}

{Per-branch patched versions where the product supports several lines, with unconfirmed ones marked.}

{**This needs action beyond upgrading.** — the log review or credential rotation, and why upgrading alone does not cover it. Omit the heading only if there is genuinely no such action.}

### Workarounds

{The mitigation, and what it does not cover.}

### References

- {fix commit or upstream PR}
- {related advisory or CVE}
- https://github.com/pgEdge/{repo}

<!-- ^^^^^^^^^^^^^^  STOP COPYING ON THE LINE ABOVE  ^^^^^^^^^^^^^^ -->

## Field 4 — Affected products

{Single-branch form:}

```
Ecosystem            {Other | Go}
Package name         {(leave empty) | module path}
Affected versions    <= {version}
Patched versions     {version}
```

{Multi-branch form — one "Add affected product" block per supported line, same repo every time:}

```
Block {N}
  Ecosystem            {Other}
  Package name         (leave empty)
  Affected versions    <= {version}
  Patched versions     {version}   {*** UNCONFIRMED — verify with WHOM ***}
```

{One line stating the ecosystem value and why it is that value — this is a field people get wrong. It is never blank.}

{Type the versions exactly as shown, with no leading `v`.}

## Field 5 — Severity

Paste this into the **Vector string** box and leave the Calculator closed — it fills the dropdown for you:

CVSS:3.1/{AV:_/AC:_/PR:_/UI:_/S:_/C:_/I:_/A:_}

The **Severity** dropdown should then read **{Band}**, and the badge above it should read **{N.N}**. If it shows anything else, the vector did not paste cleanly — fix the paste, not the dropdowns.

{Where a metric was a judgement call: name it, give the alternative score, and point at the assessment record's reasoning with an instruction not to change it without reading that first.}

## Field 6 — Weaknesses (CWE)

Type each number into **Search by CWE** and pick it from the list.

CWE-{NNN}
CWE-{NNN}

## Field 7 — Credits

{The reporter exactly as they want to be named, with how and when they reported. | Leave empty — the finding was internal, so there is no external reporter to credit.}

---

## Before you touch the form

{Numbered preconditions, each specific and checkable. Delete this section only if nothing is outstanding.}

1. {Acknowledgement sent?}
2. {Disclosure date confirmed back to the reporter?}
3. {Patched version(s) real or committed to?}
4. {Whoever gates publication for this product signed off?}

Once those are clear:

1. Fill Fields 1–7 above, **CVE identifier → "Request CVE ID later"**.
2. **Create draft security advisory.** Nothing is public: only people with access to this repository's Security tab can see it.
3. On the saved draft, press **Request CVE**.
4. **Do not press Publish advisory** until {the condition}, and every patched version named in Field 4 is confirmed to exist as a real tag.

---

## Filing notes — internal, do not paste

### {Scoping}

{What this advisory covers and what it deliberately excludes.}

### {Vulnerable code, independently verified}

{What was confirmed by direct inspection, at which tags, on what date.}

### {Reproduction status}

{What pgEdge reproduced, and what rests on the reporter's demonstration.}

### After the identifier arrives

Replace `Identifier: NOT YET REQUESTED` in `assessment-record.md`, and `CVE-{YYYY}-NNNNN` in `cve-record.json`'s `.cveMetadata.cveId` — both in this same folder.
