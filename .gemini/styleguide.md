# Code Review Style Guide

## Tone

- No subjective praise or criticism. No words like "great", "robust", "nice", "elegant", "clever".
- State facts only. Say what the code does, what it should do, and what the difference is.
- Be maximally concise. One or two sentences per comment. No filler.

## Comment Format

Prefix every comment with a severity label:

- `[critical]` — Bugs, data loss, security vulnerabilities, race conditions.
- `[major]` — Incorrect behavior, missing error handling, logic errors.
- `[minor]` — Style issues that significantly affect readability or maintainability.
- `[nit]` — Minor style preferences, naming, formatting.

## What to Review

- Correctness: bugs, logic errors, off-by-one, nil pointer dereferences.
- Security: injection, improper input validation, leaked credentials.
- Error handling: unchecked errors, missing context in error wrapping.
- Test coverage: missing tests, uncovered edge cases, insufficient assertions.
- Critical style issues: only flag style problems that meaningfully impact readability or maintainability.
- Duplication: only flag if 3+ copies of substantially identical code exist.

## What NOT to Review

- Do not comment on code that is correct and readable.
- Do not suggest refactors unless there is a concrete problem.
- Do not add comments about documentation or comment style.
- Do not praise or summarize the PR.

## Fix Suggestions

- When flagging an issue, include a concrete code suggestion when possible.
- Keep suggestions minimal — show only the changed lines, not surrounding context.

## Go Conventions

- Errors must be checked. Use `fmt.Errorf("context: %w", err)` for wrapping.
- Receiver names: short, consistent, not `this` or `self`.
- Interface names: single-method interfaces use `-er` suffix (e.g. `Reader`).
- Avoid naked returns in functions longer than a few lines.
- Use `context.Context` as the first parameter where applicable.
- Prefer table-driven tests.

## Summary Format

The PR summary (if enabled) must be:
- A single paragraph, 2-3 sentences maximum.
- State what changed and why. No opinions.
