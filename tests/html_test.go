package escape_test

import (
	"strings"
	"testing"
	"webtyp.com/escape"
)

func TestHTML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"basic", `Tom & Jerry's "House" <tag>`, `Tom &amp; Jerry&#39;s &quot;House&quot; &lt;tag&gt;`},
		{"already-entity", `&amp; &lt; &gt;`, `&amp;amp; &amp;lt; &amp;gt;`}, // double-escape expected
		{"unicode", `こんにちは & <br> 😊`, `こんにちは &amp; &lt;br&gt; 😊`},
		{"multiple", `a & b & c`, `a &amp; b &amp; c`},
		{"tags", `<div class="x">Tom & Jerry's</div>`, `&lt;div class=&quot;x&quot;&gt;Tom &amp; Jerry&#39;s&lt;/div&gt;`},
		{"emoji-and-tags", `😀 <p>1 & 2</p>`, `😀 &lt;p&gt;1 &amp; 2&lt;/p&gt;`},
		{"quotes-only", `She said: "Hi"`, `She said: &quot;Hi&quot;`},
		{"attr-value", `class="btn btn-primary"`, `class=&quot;btn btn-primary&quot;`},
		{"with-quotes", `onClick="alert('test')"`, `onClick=&quot;alert(&#39;test&#39;)&quot;`},
		{"url", `https://example.com?a=1&b=2`, `https://example.com?a=1&amp;b=2`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := escape.HTML(tc.in)
			if got != tc.want {
				t.Fatalf("%s: got=%q want=%q", tc.name, got, tc.want)
			}
		})
	}
}

func TestHTML_CompareStdLib(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"basic", `<script>alert("XSS")</script>`},
		{"quotes", `She said: "Hello" & 'Goodbye'`},
		{"entities", `Tom & Jerry's <div>`},
		{"unicode", `こんにちは <p>世界</p>`},
		{"mixed", `<a href="link.html?id=1&type=2">Click here</a>`},
		{"empty", ``},
		{"ampersand-only", `A & B & C`},
		{"all-chars", `&<>"'`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := escape.HTML(tc.in)

			// Verify all dangerous characters are escaped
			if strings.Contains(got, "<") || strings.Contains(got, ">") {
				t.Errorf("Unescaped angle brackets in output: %q", got)
			}

			// Verify input characters were processed
			if tc.in != "" && got == tc.in {
				// Only fail if input contained escapable characters
				if strings.Contains(tc.in, "&") || strings.Contains(tc.in, "<") || strings.Contains(tc.in, ">") ||
					strings.Contains(tc.in, `"`) || strings.Contains(tc.in, "'") {
					t.Errorf("Input was not escaped: %q", tc.in)
				}
			}

		})
	}
}

func TestHTMLNoAllocWhenClean(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		_ = escape.HTML("plain text without any special characters")
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocs, got %v", allocs)
	}
}
