# CVSS Judgement Calls

Most CVSS base metrics are mechanically deducible once the facts are established. Two are not: **Attack Complexity** and **Scope** regularly admit two defensible answers on the same facts, and the answer moves the score by a band or more.

This document is scoring judgement, not policy. It carries no timing tables and no CVE determination test — those live in the procedure document at `process/DISCLOSURE-PROCESS.md` in the records repository, and must be read there, fresh, every run.

Its purpose is that the reasoning below does not have to be reconstructed from scratch each time one of these comes up.

---

## Attack Complexity (AC)

### What the metric asks

Whether success depends on conditions outside the attacker's control, or on preparation specific to the target. **Not** how hard the attack was to invent, and **not** how skilled the attacker has to be — CVSS assumes a skilled attacker throughout. The question is whether, having built the exploit once, the attacker can fire it at any instance and expect it to work.

### What makes it ambiguous

The temptation is to read AC:H as "this looks difficult". A memory-corruption bug with a published proof of concept looks difficult to most readers and is often trivially repeatable in practice. Conversely, a bug that looks like a one-line request may in fact only work against one build, on one architecture, with one memory layout.

The wrong instinct is to score the exploit's *sophistication*. The right question is its *reliability against an arbitrary target*.

### The signals that point to AC:H

- **The exploit needs information it does not start with**, and has to obtain it first. The clearest case is an information leak used to defeat ASLR before the payload can be aimed. That is two stages, and the first has to succeed before the second can be built.
- **The payload is architecture-specific or build-specific.** A chain written for little-endian ARM does not run on x86-64, and one built against particular library offsets does not survive a different build of the same version.
- **A precondition exists that the attacker cannot create**, such as a configuration the operator happens to have set, or a race the attacker can only win sometimes.

### The signals that point to AC:L

- A single request, repeatable against any instance, with the same result every time.
- No target-specific preparation — no address leak, no offsets, no fingerprinting.
- Any precondition is either the default configuration or something the attacker can themselves bring about.

### The question to put to the user, in plain language

> Does this attack work the same way against every installation, or does the attacker have to learn something about the specific machine first and then build a custom payload for it? If they can write one exploit and fire it anywhere, that is Low complexity. If they have to do reconnaissance against each target and tailor the attack to what they find, that is High.

Follow it with the two scores, because the gap is usually large. Moving AC from High to Low with everything else unchanged typically moves a High into Critical, and the user should know that is the decision in front of them.

### Worked reasoning, both directions

**For AC:H, on a two-stage memory-corruption chain:** the first crafted input produces an out-of-bounds read that discloses adjacent process memory, which is what makes library and heap addresses knowable. Only then can a second input be built, carrying a forged structure at an address that is only valid for that process on that machine. The reporter's own notes naming a specific endianness and architecture are evidence for this reading. Success depends on conditions beyond the attacker's control and on specialised per-target preparation — which is what CVSS's AC:H describes.

**For AC:L, on the same shape of bug:** if the leak step is not actually required — because the target is not randomised, or the useful address is constant across builds — then the "two-stage" description is a description of how the reporter happened to write it, not of what the attack requires. A chain that can be reduced to one reliable request is AC:L regardless of how it was first demonstrated. Check whether stage one is genuinely necessary before accepting AC:H.

---

## Scope (S)

### What the metric asks

Whether the impact stays inside the same security authority as the vulnerable component, or crosses into a different one. "Authority" here means the thing that decides what is permitted — an OS user, a sandbox, a database's own permission system — not a machine and not a process.

### What makes it ambiguous

CVSS's own worked example is the source of the confusion, and it is worth knowing precisely because both readings of it turn up in real arguments.

The spec contrasts two attacks on the same database. A SQL injection that reads and writes data the database already governs is **Scope: Unchanged** — the vulnerable component and the impacted component are both the database, under the same authority. A SQL injection that reaches `xp_cmdshell` and executes OS commands is **Scope: Changed** — the impact has left the database's authority and landed in the operating system's.

The hard cases sit between those. A memory-corruption bug inside a database server process that yields code execution is genuinely arguable, because the database server process *already runs as* the OS user the attacker ends up as. Nothing moved between authorities in the way the `xp_cmdshell` example moves; and yet the attacker started with a database privilege and finished with an OS one.

### The two readings

**Unchanged.** The vulnerable component is the server process itself, which was already running under that OS user. Executing arbitrary code inside a process that already held that authority is not a crossing — it is misuse of the authority the component always had. This is broadly how memory-corruption RCEs in database servers and their extensions have been scored in practice, which is worth weighing for consistency.

**Changed.** Look at what the attacker held rather than what the component held. They began with an ordinary database privilege — no superuser, nothing resembling permission to run OS commands — and ended with full OS execution. That is precisely the boundary that hardening like restricting `COPY ... TO PROGRAM` to superusers exists to enforce, and this bug routes around that check entirely. If a boundary exists and the attack crosses it, that is what Changed is for.

Both are defensible. The choice is the user's.

### The question to put to the user, in plain language

Do not open with the metric name. Open with what the attacker gains:

> The attacker starts out as an ordinary database user with one specific permission, and ends up able to run any command on the server as the operating-system account the database runs under. The scoring question is whether that counts as breaking out of the database into the operating system, or whether it does not, because the database server process was already running as that account — so nothing new was technically entered, only misused.
>
> Scoring it as a break-out gives *X*. Scoring it as staying inside gives *Y*. Both are defensible and other vendors have gone both ways on this exact shape of bug.

### If that does not land

Drop the metric entirely and use an analogy built from the vulnerability's own facts. Something of this shape:

> Think of it as a building. A member of staff with a key to one office finds a flaw that lets them do anything the building's caretaker can do. One reading says they have broken out of their office into the caretaker's authority, and the score goes up. The other says the flaw was *in the caretaker's own toolkit*, which already had that power — so nothing was broken out of, the toolkit was just misused, and the score stays where it is.
>
> Which of those two descriptions do you think fits better here?

Then give the two numbers again, plainly. Do not re-ask the original question.

### Recording it

Whichever way it goes, write **both cases** into the assessment record along with the alternative score, and make clear it was a deliberate decision rather than a default. Note in the advisory draft that the Scope value was chosen deliberately and should not be changed without re-reading the record's reasoning — otherwise someone will "fix" it at filing time.

---

## A Note on the Metrics That Are Not Judgement Calls

Do not manufacture a question where there is not one. Asking the user about Privileges Required when the answer follows directly from the verified facts spends the attention you need for the one metric that genuinely needs it.

But do **record** the reasoning for a metric that is quietly deciding the outcome. Where Attack Complexity being High is the only thing holding a vulnerability out of the Critical band, that sentence belongs in the record with the alternative score next to it — stated as settled, not offered as a question. It is the sentence someone will challenge later, and it should already be answered when they do.
