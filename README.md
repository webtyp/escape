# escape
<img src="docs/img/badges.svg">

Make text safe for an output context: HTML, HTML attributes, JSON strings.

## Context mapping

| I want... | Use |
|---|---|
| HTML text or attribute | `escape.HTML(s)` |
| JSON string | `escape.JSON(b, s)` |

### Migration from `webtyp.com/fmt`

| Old (`webtyp.com/fmt`) | New (`webtyp.com/escape`) |
|---|---|
| `fmt.Convert(s).EscapeHTML()` | `escape.HTML(s)` |
| `fmt.Convert(s).EscapeAttr()` | `escape.HTML(s)` |
| `fmt.JSONEscape(s, b)` | `escape.JSON(b, s)` |
| `fmt.Html(format, args...)` | `fmt.Sprintf(format, args...)` |

## Warning

`escape.HTML(s)` does **NOT** make text safe for URLs, JS or CSS contexts. It only protects against injections within standard HTML text elements or within HTML attributes quoted with single or double quotes.
