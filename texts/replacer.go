package texts

import (
	"regexp"
	"strings"
)

type TextReplacer interface {
	Execute(string) string
}

type Replacer struct {
	replacers []TextReplacer
}

func NewReplacer() *Replacer {
	return &Replacer{}
}

func (r *Replacer) AddRegexp(src *regexp.Regexp, replaced string) {
	r.replacers = append(r.replacers, NewRegexpReplacer(src, replaced))
}

func (r *Replacer) AddString(src string, replaced string) {
	r.replacers = append(r.replacers, NewStringReplacer(src, replaced))
}

func (r *Replacer) AddFunc(f func(string) string) {
	r.replacers = append(r.replacers, FuncReplacer(f))
}

func (r *Replacer) Execute(value string) string {
	for _, expr := range r.replacers {
		value = expr.Execute(value)
	}

	return value
}

type RegexpReplacer struct {
	src      *regexp.Regexp
	replaced string
}

func NewRegexpReplacer(src *regexp.Regexp, replaced string) *RegexpReplacer {
	return &RegexpReplacer{src: src, replaced: replaced}
}

func (r *RegexpReplacer) Execute(value string) string {
	return r.src.ReplaceAllString(value, r.replaced)
}

type StringReplacer struct {
	src      string
	replaced string
}

func NewStringReplacer(src string, replaced string) *StringReplacer {
	return &StringReplacer{src: src, replaced: replaced}
}

func (r *StringReplacer) Execute(value string) string {
	return strings.ReplaceAll(value, r.src, r.replaced)
}

type FuncReplacer func(string) string

func (f FuncReplacer) Execute(value string) string {
	return f(value)
}
