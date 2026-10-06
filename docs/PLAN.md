---
PLAN: "feat: escape — HTML and JSON escaping moved out of webtyp.com/fmt"
TAG: v0.1.0
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# Plan — escape: make text safe for an output context

Phase **A5 (gate)** of the master plan `SOURCE_SELECTION_MASTER_PLAN.md` (orchestration only —
everything this plan needs is inline). `webtyp/dom`, `webtyp/json`, `veltylabs/sitetheme` and
`webtyp/goflare-demo` migrate after this tag exists.

Read [AGENTS.md](../AGENTS.md) first (WASM rules, no state/`init`, tests in `tests/`).

## Why

Escaping lived in `fmt` as `Conv` methods (`EscapeHTML`, `EscapeAttr`) plus `JSONEscape`. It is a
separate responsibility. Its natural consumer, `webtyp/dom`, cannot be its home, for two reasons:
- `webtyp/html` depends on `dom`, so putting it in `html` would be an import cycle.
- In WASM, `dom` runs code at package load (`instance = newDom(shared)` in `dom/dom.go`). That reads
  `localStorage`/`document` and keeps the whole DOM implementation in the binary. A server or Worker
  that only escapes an email body would pay for a DOM engine.

A dependency-light package under `dom` serves every consumer.

## State when this plan starts (moved by the maintainer, 2026-10-06)

Moved **unchanged** (they still say `package fmt` and use `Conv` methods):
- `html.go` = `webtyp/fmt/html.go`: `(*Conv).EscapeAttr`, `(*Conv).EscapeHTML`, `Html(...)`.
- `json.go` = the `JSONEscape(s string, b *Builder)` function cut from `webtyp/fmt/quote.go`.
- `tests/html_test.go` = `webtyp/fmt/html_test.go`: `TestHtml`, `BenchmarkHtml`, `TestEscapeAttr`,
  `TestEscapeHTML_TableDriven`, `TestEscapeHTML_CompareStdLib`, `TestEscapeAttr_CompareStdLib`.
- `tests/json_test.go` = `webtyp/fmt/json_escape_test.go`.
- `docs/API_HTML.md`, `docs/API_JSON_ESCAPE.md`.

The gonew stub was deleted. Do not recreate it.

## Design gate

1. **Prior art.** Go's `html.EscapeString` and `template.HTMLEscapeString` (package `html`,
   separate from `fmt`); `encoding/json`'s internal string escaping; Rust's `html-escape` and
   `v_htmlescape` crates; OWASP's per-context encoder (HTML body, attribute, JS, URL). Everyone keeps
   escaping out of the formatter and names it by output context. That is what this repo does.
2. **Novice-name test.** `escape.HTML(s)` — "escape for HTML". `escape.JSON(b, s)` — "escape s as
   JSON into b". Argument order is destination then source, as in `copy(dst, src)`.

   | Old (`webtyp.com/fmt`) | New (`webtyp.com/escape`) |
   |---|---|
   | `fmt.Convert(s).EscapeHTML()` | `escape.HTML(s)` |
   | `fmt.Convert(s).EscapeAttr()` | `escape.HTML(s)` |
   | `fmt.JSONEscape(s, b)` | `escape.JSON(b, s)` |
   | `fmt.Html(format, args...)` | `fmt.Sprintf(format, args...)` (both real callers use format mode, without `lang`) |
3. **Complexity ledger.**
   ```
   Concepts the developer must learn   +0 / −2 (EscapeAttr, Html)
   Files they must touch to do X       +0 / −0
   Lines at the call site              ±0
   Ways to do the same thing           −2  (EscapeAttr ≡ EscapeHTML: both escape exactly & < > " ' ; Html ≡ Sprintf)
   ```
   `EscapeHTML` and `EscapeAttr` escaped the **same five characters to the same entities**: only the
   order of the `Replace` calls differed. They are one function.
4. **Where it belongs.** Its own repo below `dom` (see Why).
5. **What it deletes.** `EscapeAttr`, `Html`, `TestHtml`, `BenchmarkHtml`, and the
   `Conv`-method form.

## Stage 1 — `html.go` (rewrite)

```go
package escape

// HTML returns s with & < > " ' replaced by &amp; &lt; &gt; &quot; &#39;.
// The result is safe as HTML text and inside a quoted attribute value
// (single or double quotes). It does NOT make text safe for URLs, JS or CSS.
func HTML(s string) string
```

One pass: scan for the first byte in `&<>"'`; if none, `return s` (zero allocations). Otherwise
`buf := make([]byte, 0, len(s)+16)`, copy the clean prefix, then append bytes or entities, and return
`string(buf)`. Delete `EscapeAttr`, `EscapeHTML` and `Html`.

## Stage 2 — `json.go`

```go
package escape

import "webtyp.com/fmt"

// JSON writes s into b escaped for the inside of a JSON string (no surrounding
// quotes): " → \"  \ → \\  newline → \n  CR → \r  tab → \t  other bytes < 0x20 → \u00XX.
// The caller writes the quotes, so strings compose without extra allocations.
func JSON(b *fmt.Builder, s string)
```

Body: the moved `JSONEscape` loop unchanged; only the signature changes (`b` first). `go get
webtyp.com/fmt@latest`.

## Stage 3 — tests in `tests/`

`package escape_test`, importing `webtyp.com/escape` (and `webtyp.com/fmt` for `Builder`/`Contains`).
- Delete `TestHtml` and `BenchmarkHtml`.
- `Convert(x).EscapeHTML()` and `Convert(x).EscapeAttr()` → `escape.HTML(x)`. Keep every expected
  output (both old functions produced the same output).
- Merge `TestEscapeAttr` and `TestEscapeAttr_CompareStdLib` cases into the HTML tests (dedupe).
- `JSONEscape(s, b)` → `escape.JSON(b, s)`.
- Add `TestHTMLNoAllocWhenClean`: `testing.AllocsPerRun(100, func() { _ = escape.HTML("plain text") })`
  must be `0`.

## Stage 4 — docs

`README.md`: an "I want X → use Y" table (HTML text or attribute → `HTML`; JSON string → `JSON`), the
old→new table above, and the warning that `HTML` does not cover URL/JS/CSS contexts. Then delete
`docs/API_HTML.md` and `docs/API_JSON_ESCAPE.md`.

## Acceptance

- `gotest` passes (includes WASM).
- `ls *_test.go 2>/dev/null` → nothing at the root.
- `grep -rn "EscapeAttr\|EscapeHTML\|JSONEscape\|func Html\|Conv)" --include='*.go' .` → empty.

## Stages

| # | Stage | Files |
|---|---|---|
| 1 | HTML | `html.go` |
| 2 | JSON | `json.go`, `go.mod`, `go.sum` |
| 3 | Tests | `tests/html_test.go`, `tests/json_test.go` |
| 4 | Docs | `README.md`, `docs/*.md` (deleted) |
