# Phase 2 — CVSS Scoring

The score decides the severity band, which decides the notification clock and who has to be told before publication. Getting it wrong in either direction is expensive, so it is built metric by metric with reasoning attached, and the metrics that are genuinely arguable are put to the user rather than guessed.

Use CVSS v3.1 unless the live procedure document says otherwise.

## Step 1: Build the Vector One Metric at a Time

Work through all eight base metrics in order — AV, AC, PR, UI, S, C, I, A — and for each one write a sentence of reasoning grounded in the facts established in Phase 1. Not "PR:L because low privileges", but "PR:L because the attacker needs `CREATE` on a schema that is a replication-set member, not superuser".

Most metrics are mechanically deducible from the verified facts. State them with their reasoning and move on. Do not ask the user about a metric that has one defensible answer — that wastes their attention on the metrics where it is not needed, which is exactly the attention you want available for the one where it is.

The impact metrics (C, I, A) follow from the demonstrated outcome. If the outcome is arbitrary code execution as the service's OS user, all three are High and there is nothing to debate.

## Step 2: Escalate the Genuine Judgement Calls

A metric needs the user's decision when both answers are defensible on the same facts and the choice moves the score. Do not resolve it silently, and do not resolve it by picking the higher number "to be safe" — an inflated score is as much a defect as a deflated one.

`standards/CVSS-JUDGMENT-CALLS.md` carries the worked reasoning for the two metrics that most often land here, Attack Complexity and Scope. Read it before escalating, so the question you ask is already informed.

When you escalate, use `AskUserQuestion` and put four things in front of the user:

1. **What the metric is actually asking**, in plain language, without CVSS jargon.
2. **The case for each answer**, in terms of this vulnerability's specific facts — not the general definition.
3. **The concrete score each answer produces**, and the severity band each lands in. "7.5 High versus 8.5 High" is a different decision from "7.5 High versus 9.1 Critical", and the user needs to know which one they are making.
4. **Which way comparable advisories have been scored**, if a precedent exists among published records in the records repository's `records/` directory.

### If the user says they do not understand the question

**Re-explain it differently. Do not re-ask it the same way.** A question that did not land will not land better on a second reading.

Drop the metric name entirely, and describe the thing being decided with a concrete analogy drawn from the vulnerability itself — who ends up able to do what, and whether that counts as crossing a line. Then restate the two answers as two plain-English sentences the user can simply pick between, and give the two numbers again.

Getting the user to a genuine decision is the point. A decision made by a confused user is worth no more than one you made yourself.

## Step 3: Record the Reasoning, Not Just the Answer

Whatever the user chooses, write down **both cases** — the one taken and the one not taken — in the assessment record. A future reader who disagrees with the score needs to see that the alternative was considered and why it lost, and a reviewer months later needs to be able to re-open the question without reconstructing it from nothing.

Also note in the record which metrics were *not* judgement calls, where one of them is doing unusual work in the final score. If Attack Complexity being High is the only thing keeping a vulnerability out of the Critical band, say so plainly, with the score the alternative would have produced. That is the sentence someone will challenge, so it should already be answered.

## Step 4: Verify the Arithmetic

Compute the base score by hand from the vector and check it against the severity band boundaries. Then write into the record that the score is hand-computed and **still needs re-confirming against GitHub's own calculator at filing time**.

This is not redundant. The filing guide's procedure is to paste the vector string into the form's vector box and check that the number GitHub then displays matches the record. If it does not, the vector did not paste cleanly — and the fix is to correct the paste, never to adjust the dropdowns until the number agrees.

## Step 5: Choose CWEs

Pick **two or three**, following the pattern already established in this program's filed advisories: root cause, consequence, mechanism.

- **Root cause** — the thing that was done wrong. Improper validation of a specified quantity, improper access control, improper neutralisation.
- **Consequence** — what it produces. Out-of-bounds read, SQL injection, information exposure.
- **Mechanism** — how the consequence becomes an impact, where that is a distinct step. Untrusted pointer dereference, for instance.

Where two of the three collapse into one CWE, use two rather than padding to three. Check the CWEs used by comparable records in `records/*/cve-record.json` before choosing — consistency across the program's advisories is worth more than a marginally more precise identifier.

## Output of This Phase

- the full vector string, in the form `CVSS:3.1/AV:.../A:...`
- the base score and severity band, hand-computed
- a sentence of reasoning per metric
- for any escalated metric: both cases, the user's decision, the date, and the alternative score
- two or three CWEs with their full names

Proceed to Phase 3.
