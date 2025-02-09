package texts

import "strings"

// SimplePattern is a simple pattern for matching strings.
//
// It supports the following special characters:
// - `* - matches any sequence of characters.
// - `$` - matches the end of the string.
type SimplePattern string

func (p SimplePattern) MinLength() int {
	constants := strings.ReplaceAll(string(p), "*", "")
	constants = strings.ReplaceAll(constants, "$", "")

	return len(constants)
}

func (p SimplePattern) IsSuffix() bool {
	return strings.HasSuffix(string(p), "$")
}

func (p SimplePattern) IsGreedy() bool {
	return strings.HasSuffix(string(p), "*")
}

func (p SimplePattern) Parts() []string {
	pattern := string(p)
	pattern = strings.ReplaceAll(pattern, "$", "")
	parts := strings.Split(pattern, "*")

	return parts
}

func (p SimplePattern) String() string {
	return string(p)
}
