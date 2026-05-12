---
name: feedback-loop
description: Internal subagent for iterating on a GitHub PR (or non-PR branch) until CI is green and review threads are resolved. Dispatched by pgedge-fix-issue, pgedge-monitor-actions, and pgedge-review-pr. Not directly user-invoked - users invoke the caller skills.
allowed-tools: Bash(gh:*), Bash(git:*), Bash(make:*), Bash(npm:*), Bash(pytest:*), Bash(go:*), Bash(cargo:*), Bash(ruff:*), Bash(sleep:*), Read, Edit, Write, Grep
---

# Feedback-Loop Primitive

<!-- BODY-START -->
<!-- Sections inserted by subsequent tasks. -->
<!-- BODY-END -->
