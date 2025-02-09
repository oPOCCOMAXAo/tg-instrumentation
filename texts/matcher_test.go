package texts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatcher_Match(t *testing.T) {
	testCases := []struct {
		value   string
		pattern SimplePattern
		result  int
		err     string
	}{
		{
			value:   "",
			pattern: "$$",
			err:     "more than one $ in pattern: $$",
		},
		{
			value:   "/",
			pattern: "/",
			result:  1,
		},
		{
			value:   "/a/b",
			pattern: "/*",
			result:  4,
		},
		{
			value:   "/a/b",
			pattern: "/",
			result:  1,
		},
		{
			value:   "/a/b",
			pattern: "/*/b",
			result:  4,
		},
		{
			value:   "/a/b",
			pattern: "/a/b",
			result:  4,
		},
		{
			value:   "/a/b",
			pattern: "/a/bc",
			result:  -1,
		},
		{
			value:   "/a/b",
			pattern: "*",
			result:  4,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.pattern.String(), func(t *testing.T) {
			matcher, err := NewSimpleMatcher(tC.pattern)
			if tC.err != "" {
				require.ErrorContains(t, err, tC.err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tC.result, matcher.Match(tC.value))
		})
	}
}
