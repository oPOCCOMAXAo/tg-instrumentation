package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		input  string
		output *Query
	}{
		{
			input: "",
			output: &Query{
				Params: map[string][]string{},
			},
		},
		{
			input: "/add",
			output: &Query{
				Command: "/add",
				Params: map[string][]string{
					"/add": nil,
				},
			},
		},
		{
			input: "multi word=2,3 params=4",
			output: &Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": {"4"},
				},
			},
		},
		{
			input: "multi word=2 word=3 params",
			output: &Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": nil,
				},
			},
		},
		{
			input: "command=test with params",
			output: &Query{
				Command: "command",
				Params: map[string][]string{
					"command": {"test"},
					"with":    nil,
					"params":  nil,
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			res := Decode(tC.input)
			require.Equal(t, tC.output, res)
		})
	}
}

func TestQuery_Encode(t *testing.T) {
	testCases := []struct {
		input  Query
		output string
	}{
		{
			input: Query{
				Params: map[string][]string{},
			},
			output: "",
		},
		{
			input: Query{
				Params: map[string][]string{
					"add": nil,
				},
			},
			output: "add",
		},
		{
			input: Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": {"4"},
				},
			},
			output: "multi params=4 word=2,3",
		},
		{
			input: Query{
				Command: "command",
				Params:  map[string][]string{},
			},
			output: "command",
		},
		{
			input: Query{
				Command: "command",
				Params: map[string][]string{
					"command": {"test"},
					"with":    nil,
					"params":  nil,
				},
			},
			output: "command=test params with",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.output, func(t *testing.T) {
			res := tC.input.Encode()
			require.Equal(t, tC.output, res)
		})
	}
}

func TestQuery_setters(t *testing.T) {
	res := Command("cmd").
		WithParam("string", "s").
		WithParamInt64("int", 42).
		WithParamEmpty("empty").
		Encode()

	require.Equal(t, "cmd empty int=42 string=s", res)
}

func TestQuery_getters(t *testing.T) {
	query := Decode("cmd empty int=42 string=s")

	{
		require.True(t, query.Has("empty"))

		str, ok := query.Get("empty")
		require.True(t, ok)
		require.Empty(t, str)

		i, ok := query.GetInt64("empty")
		require.True(t, ok)
		require.Empty(t, i)
	}

	{
		var value string

		str, ok := query.Get("string")
		require.True(t, ok)
		require.Equal(t, "s", str)

		require.True(t, query.GetInto("string", &value))
		require.Equal(t, "s", value)
	}

	{
		var value int64

		str, ok := query.Get("int")
		require.True(t, ok)
		require.Equal(t, "42", str)

		ok = query.GetInt64Into("int", &value)
		require.True(t, ok)
		require.Equal(t, int64(42), value)
	}
}
