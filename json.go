package escape

import "webtyp.com/fmt"

// JSON writes s into b escaped for the inside of a JSON string (no surrounding
// quotes): " → \"  \ → \\  newline → \n  CR → \r  tab → \t  other bytes < 0x20 → \u00XX.
// The caller writes the quotes, so strings compose without extra allocations.
func JSON(b *fmt.Builder, s string) {
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue // safe byte — part of the current bulk run
		}
		if i > start {
			b.WriteString(s[start:i])
		}
		switch c {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteString(`\u00`)
			_ = b.WriteByte("0123456789abcdef"[c>>4])
			_ = b.WriteByte("0123456789abcdef"[c&0xf])
		}
		start = i + 1
	}
	if start < len(s) {
		b.WriteString(s[start:])
	}
}
