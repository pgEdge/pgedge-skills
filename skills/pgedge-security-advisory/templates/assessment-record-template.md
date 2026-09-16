# AR-{YYYY}-{NNN} {SHORT TITLE}

{Optional opening note: state here if this record breaks new ground — first external report, first multi-branch product, first advisory for this product — so the next reader knows there was no prior pattern to match against. Delete if there is a precedent.}

```
Reference:          AR-{YYYY}-{NNN}   Product: {PRODUCT}
                                  Category: {Cloud Service | Distributed Software}
Reported:           {YYYY-MM-DD}    Source: {security@ (external — REPORTER) |
                     {HH:MM}         internal | scan | pentest | public}
Acknowledged:       {date | n/a, internal | PENDING — due {DATE}, being
                     {N} business days: {list the dates explicitly}}
Reporter:           {NAME | pgEdge engineering | anonymous}
                     Credit preference: {how they want to be named, or
                     "none, internal". Note any second party credited
                     upstream.}

REACHABLE AS SHIPPED
Confirmed:                        {yes/no + evidence. Separate what pgEdge
                                  independently verified from what rests on
                                  the reporter's demonstration.}
Pinned by test:                   {test file, or n/a}

AFFECTED VERSIONS            (verified against release tags, not HEAD alone)
{Single-branch form:}
Earliest / latest affected:       {tag} / {tag}
{Multi-branch form — one stanza per concurrently-supported line:}
  Major {N}.x  first released {DATE}, supported to {DATE}
               earliest {tag}, latest GA {tag} — AFFECTED

Vulnerable code evidence ({verified DATE by direct inspection of REPO}):
  {tag}:                  {path:line}
  {tag}:                  {path:line}
{One sentence naming the pattern itself, in your own words, and stating it
is identical across those tags.}

Fixed in:                         {version | UNRELEASED + state of the fix.
                                  Mark any version named ahead of release as
                                  provisional, and name who confirms it.}

SEVERITY
CVSS v3.1 vector:                 {AV:_/AC:_/PR:_/UI:_/S:_/C:_/I:_/A:_}
Base score:                       {N.N} {Band}   Verified against vector: yes
                                  (hand-computed; confirm against GitHub's
                                  own calculator when filing)
Re-rated from CVSS:               {no | yes + rationale}

  {Judgement call, if any: which metric, who decided, on what date.}
  {The case for the value taken.}
  {The case for the value not taken, and the score it would have given.}
  {Any metric quietly deciding the band, stated as settled with the
  alternative score.}

CVE
Released-version test:            {pass/fail + evidence}
Trust-boundary test:              {pass/fail + name the untrusted party and
                                  the privilege they hold}
Outcome:                          {CVE | no CVE + reason}
Identifier:                       NOT YET REQUESTED. Do not request until:
                                  {numbered preconditions}

DISCLOSURE
Default is private until release: {followed | deviated + reasoning + date}
Reporter told outcome:            {date | PENDING — due {DATE} | n/a}
Publication date agreed:          {date, and whether it is CONFIRMED back to
                                  the reporter or merely the policy default}

NOTIFIED
Class 1 (individual, evidence):   {customer, date, evidence link | NONE as at
                                  {DATE} | NOT YET ASSESSED for this product}
Notification text:                {link | n/a}

PUBLISHED
Release:                          {version, date | PENDING}
Advisory:                         {link, date | PENDING}
                                  CVE in release notes: {yes/no/PENDING}

Signed off by:                    {name, date | PENDING}
```

## {Scoping — what this record covers and excludes}

{Where the engineering fix addresses more than was reported: state what is
covered, what is not, and why. Note that the excluded items likely want
their own record, so they are not lost.}

## {Open process questions}

{Gaps flagged rather than resolved: does the existing publish-approval chain
cover this product; does any customer agreement carry a notification
obligation here; does existing tooling handle this record's shape.}

## Customer action beyond upgrading

{Log review, credential rotation, treating old output as sensitive.
Upgrading stops further exploitation; it does not say whether the window was
already used, and readers will not infer the step. Omit only if there is
genuinely none.}

## Workaround {(until a patched version ships)}

{The mitigation, and honestly what it does not cover.}
