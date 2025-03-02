package texts

import (
	"regexp"
	"strings"
)

//nolint:gochecknoglobals
var htmlReplacer = NewHTMLReplacer()

func NewHTMLReplacer() TextReplacer {
	res := NewReplacer()
	res.AddString(`&nbsp;`, " ")
	res.AddString(`&lt;`, "<")
	res.AddString(`&gt;`, ">")
	res.AddString(`&quot;`, `"`)
	res.AddString(`&apos;`, `'`)
	res.AddString(`&amp;`, "&")
	res.AddString(`&#039;`, "'")
	res.AddString(`<br>`, "\n")
	res.AddString(`<hr>`, "\n")
	res.AddRegexp(regexp.MustCompile(`<[^<>]+>`), " ")
	res.AddRegexp(regexp.MustCompile(`[\t ]+`), " ")

	// Trim beginning and ending spaces
	res.AddRegexp(regexp.MustCompile(`(?m)^\s+|\s+$`), "") // in multiple lines
	res.AddFunc(strings.TrimSpace)                         // final trim

	return res
}

func EscapeHTML(s string) string {
	return htmlReplacer.Execute(s)
}
