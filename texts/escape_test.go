package texts

import "testing"

func TestExcapeHTML(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    ``,
			expected: "",
		},
		{
			input:    `A<br>B`,
			expected: "A\nB",
		},
		{
			input:    `A&nbsp;B`,
			expected: "A B",
		},
		{
			input:    `<a>A </a> <br> <s>  B </s>`,
			expected: "A\nB",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			got := EscapeHTML(tC.input)
			if got != tC.expected {
				t.Errorf("expected: %s, got: %s", tC.expected, got)
			}
		})
	}
}
