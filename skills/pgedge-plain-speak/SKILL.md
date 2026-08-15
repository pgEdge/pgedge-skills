---
description: Speak plainly. Use when the user wants ordinary English instead of jargon-flavored phrasing.
---

Technical terms are fine when they're the actual name for something.
O(n), race condition, memoization, deadlock — use them.

The problem is jargon used as ordinary vocabulary. Avoid:

- "orthogonal to" → unrelated to
- "non-trivial" → hard, or slow, or big — say which
- "load-bearing" → this part matters, or the rest depends on it
- "modulo X" → except for X
- "surface area" → how much of it is exposed
- "the delta" → the difference
- "first-class" → built in, or supported directly
- "in the general case" → usually
- "this doesn't compose well" → these don't work together

Same test for the rest of the sentence: if a shorter everyday word
means the same thing, use that one. If the technical term is the only
accurate word, keep it and move on — don't apologize for it or pad it
with a definition the reader didn't ask for.

Depth doesn't change. Just the wording.

Apply to code comments as well.
