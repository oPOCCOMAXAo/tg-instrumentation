package texts

import (
	"strings"

	"github.com/pkg/errors"
)

// SimpleMatcher is a simple pattern matcher.
//
// It supports the following special characters:
// - `* - matches any sequence of characters.
// - `$` - matches the end of the string.
type SimpleMatcher struct {
	parts  []string
	suffix bool
	greedy bool
	minLen int
}

func NewSimpleMatcherFromString(pattern string) (*SimpleMatcher, error) {
	return NewSimpleMatcher(SimplePattern(pattern))
}

func NewSimpleMatcher(pattern SimplePattern) (*SimpleMatcher, error) {
	var res SimpleMatcher

	res.minLen = pattern.MinLength()
	res.suffix = pattern.IsSuffix()
	res.greedy = pattern.IsGreedy()

	if strings.Count(string(pattern), "$") > 1 {
		return nil, errors.Wrapf(ErrInvalidPattern, "more than one $ in pattern: %s", pattern)
	}

	res.parts = pattern.Parts()

	if len(res.parts) == 0 {
		return nil, errors.Wrap(ErrFailed, pattern.String())
	}

	return &res, nil
}

// Match checks if the value matches the pattern.
//
// Returns the length of the matched prefix or -1 if the value does not match the pattern.
func (m *SimpleMatcher) Match(value string) int {
	if len(value) < m.minLen {
		return -1
	}

	// First part is prefix.
	index := strings.Index(value, m.parts[0])
	if index != 0 {
		return -1
	}

	matched := len(m.parts[0])
	if len(m.parts) == 1 {
		return matched
	}

	for _, part := range m.parts[1:] {
		idx := strings.Index(value[matched:], part)
		if idx == -1 {
			return -1
		}

		matched += idx + len(part)
	}

	if m.suffix {
		if !strings.HasSuffix(value, m.parts[len(m.parts)-1]) {
			return -1
		}

		return len(value)
	}

	if m.greedy {
		return len(value)
	}

	return matched
}
