# Phase 2 — Report

Generate two outputs from the audit results: a terminal summary
and a persistent file report.

## Terminal Summary

Print a compact summary using this format:

```
pgEdge Repo Ready Audit — <PROJECT_NAME>
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Detected: <LANGUAGES>
Project type: <PROJECT_TYPE>
Postgres: <yes/no>

Tier 1 — Governance           X/Y  <bar>  NN%
Tier 2 — Code Quality         X/Y  <bar>  NN%
Tier 3 — Documentation        X/Y  <bar>  NN%
Tier 4 — Developer Experience X/Y  <bar>  NN%
Tier 5 — Go                   X/Y  <bar>  NN%  (if applicable)
Tier 5 — Python               X/Y  <bar>  NN%  (if applicable)
Tier 5 — TypeScript           X/Y  <bar>  NN%  (if applicable)
Tier 6 — Postgres             X/Y  <bar>  NN%  (if applicable)

Overall: XX/YY (NN%)
```

For the progress bar, use filled and empty blocks proportional to
the percentage. 14 characters wide.

Example: `5/7  ██████████░░░░  71%`

After the summary, print the remediation prompt:

```
Remediation options:
  [A] Apply all fixes
  [P] Pick categories to fix
  [O] Walk through one at a time
  [S] Skip remediation (audit only)

Choose [A/P/O/S]:
```

### Re-run Indicators

When running on a previously onboarded repo (REPO_READY_AUDIT.md
already exists), add status indicators:

- Checks that were PASS and are still PASS: show normally
- Checks that were PASS but are now not: mark as **REGRESSION**
- Checks not in the previous audit (new standard): mark as **NEW**

Show regression count in the summary header if > 0:

```
Regressions detected: 2
```

## File Report — REPO_READY_AUDIT.md

Write the detailed report to `REPO_READY_AUDIT.md` at the repo
root.

### Ensure .gitignore Coverage

Before writing, check that `.gitignore` contains
`REPO_READY_AUDIT.md`. If not, append it:

```bash
grep -q "REPO_READY_AUDIT.md" .gitignore 2>/dev/null \
  || echo "REPO_READY_AUDIT.md" >> .gitignore
```

### Report Format

```markdown
# pgEdge Repo Ready Audit Report

**Project:** <PROJECT_NAME>
**Date:** <YYYY-MM-DD>
**Languages:** <LANGUAGES>
**Project type:** <PROJECT_TYPE>
**Overall score:** XX/YY (NN%)

---

## Tier 1 — Governance (X/Y)

### G-01: License — <STATUS>
**Found:** <what was observed>
**Expected:** LICENSE.md with PostgreSQL License, pgEdge copyright
**Remediation:** Copy templates/license/LICENSE.md, substitute
placeholders

### G-02: CODEOWNERS — <STATUS>
...

---

## Tier 2 — Code Quality (X/Y)
...

(repeat for all tiers and checks)
```

### Status Formatting

- **PASS** — no remediation section needed
- **PARTIAL** — include what's present and what's missing
- **MISSING** — include remediation action
- **SKIP** — one-line note on why (e.g., "Go not detected")
- **REGRESSION** — bold, include what changed from last audit
- **NEW** — note that this is a newly added standard

## Transition to Phase 3

After the user chooses a remediation mode:

- **A (Apply all):** Proceed to Phase 3 with mode=batch
- **P (Pick categories):** Ask which tiers to remediate, then
  proceed to Phase 3 with mode=selected and the chosen tiers
- **O (One at a time):** Proceed to Phase 3 with mode=interactive
- **S (Skip):** End the skill run. Print "Audit complete. Report
  saved to REPO_READY_AUDIT.md."
