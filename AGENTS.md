# AGENTS.md — webtyp/escape

Working notes for AI agents operating in this repository. End-user docs: [README.md](README.md).

## What this repo is

Make text safe for an output context: `escape.HTML(s)` for HTML text and quoted attribute values, and
`escape.JSON(b, s)` for the inside of a JSON string. It has no state and no `init`, so anything can
import it — `webtyp/dom` (SSR serialization), `webtyp/json`, servers and email bodies — without
pulling in a DOM engine.

## This package compiles to WASM

- Do NOT import `strings`, `strconv`, `errors`, stdlib `fmt` or stdlib `html`. The only dependency is
  `webtyp.com/fmt` (for `fmt.Builder`).
- No `map`, no `reflect`, no `init`, no package-level state.
- Fast path first: input with nothing to escape is returned as-is, with zero allocations.

## The build that defines "done"

```bash
go install webtyp.com/devflow/cmd/gotest@latest   # once
gotest
```

## Rules

- Tests live in `tests/` as `package escape_test` (public API only). A root-level test is allowed only
  with a top-of-file justification of the unexported identifier it needs. **Never export a symbol so
  a test can reach it.**
- `webtyp/dom` imports this package. This package must never import `dom` or `html`.
