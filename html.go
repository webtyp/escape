package escape

// HTML returns s with & < > " ' replaced by &amp; &lt; &gt; &quot; &#39;.
// The result is safe as HTML text and inside a quoted attribute value
// (single or double quotes). It does NOT make text safe for URLs, JS or CSS.
func HTML(s string) string {
	i := 0
	for i < len(s) {
		c := s[i]
		if c == '&' || c == '<' || c == '>' || c == '"' || c == '\'' {
			break
		}
		i++
	}

	if i == len(s) {
		return s
	}

	buf := make([]byte, 0, len(s)+16)
	buf = append(buf, s[:i]...)

	for j := i; j < len(s); j++ {
		c := s[j]
		switch c {
		case '&':
			buf = append(buf, "&amp;"...)
		case '<':
			buf = append(buf, "&lt;"...)
		case '>':
			buf = append(buf, "&gt;"...)
		case '"':
			buf = append(buf, "&quot;"...)
		case '\'':
			buf = append(buf, "&#39;"...)
		default:
			buf = append(buf, c)
		}
	}

	return string(buf)
}
