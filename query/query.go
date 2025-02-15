package query

import (
	"slices"
	"strconv"
	"strings"
)

const (
	QueryParamDelimiter = " "
	QueryValueDelimiter = "="
	QuerySliceDelimiter = ","
)

type Query struct {
	Command string // Command is the first parameter in the query.
	Params  map[string][]string
}

func New() *Query {
	return &Query{
		Params: make(map[string][]string),
	}
}

func Command(command string) *Query {
	return &Query{
		Command: command,
		Params:  make(map[string][]string),
	}
}

// Decode decodes the query string into the Query structure.
//
// Text examples:
//
//	"/add"
//	"menu page=help"
//	"menu=help"
//	"article id=1"
func Decode(query string) *Query {
	res := New()
	res.Decode(query)

	return res
}

func (q *Query) WithCommand(command string) *Query {
	q.Command = command

	return q
}

func (q *Query) WithParam(key string, value string) *Query {
	q.Params[key] = append(q.Params[key], value)

	return q
}

func (q *Query) WithParamInt64(key string, value int64) *Query {
	return q.WithParam(key, strconv.FormatInt(value, 10))
}

func (q *Query) Get(key string) (string, bool) {
	values, ok := q.Params[key]
	if !ok {
		return "", false
	}

	if len(values) == 0 {
		return "", true
	}

	return values[0], true
}

func (q *Query) GetInto(key string, into *string) bool {
	value, ok := q.Get(key)
	if !ok {
		return false
	}

	*into = value

	return true
}

func (q *Query) GetSlice(key string) ([]string, bool) {
	values, ok := q.Params[key]
	if !ok {
		return nil, false
	}

	return values, true
}

func (q *Query) GetInt64(key string) (int64, bool) {
	values, ok := q.Params[key]
	if !ok {
		return 0, false
	}

	if len(values) == 0 {
		return 0, true
	}

	res, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return 0, false
	}

	return res, true
}

func (q *Query) GetInt64Slice(key string) ([]int64, bool) {
	values, ok := q.Params[key]
	if !ok {
		return nil, false
	}

	res := make([]int64, 0, len(values))

	for _, value := range values {
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, false
		}

		res = append(res, val)
	}

	return res, true
}

func (q *Query) GetInt64Into(key string, into *int64) bool {
	value, ok := q.GetInt64(key)
	if !ok {
		return false
	}

	*into = value

	return true
}

func (q *Query) Encode() string {
	parts := make([]string, 0, 1+len(q.Params))

	if q.Command != "" {
		parts = append(parts, JoinKeyValues(q.Command, q.Params[q.Command]))
	}

	params := make([]string, 0, len(q.Params))

	for key, values := range q.Params {
		if key == q.Command {
			continue
		}

		if len(values) == 0 {
			params = append(params, key)

			continue
		}

		params = append(params, JoinKeyValues(key, values))
	}

	if len(params) > 0 {
		slices.Sort(params)

		parts = append(parts, strings.Join(params, QueryParamDelimiter))
	}

	return strings.Join(parts, QueryParamDelimiter)
}

func (q *Query) Decode(query string) {
	paramsList := strings.Split(query, QueryParamDelimiter)

	q.Params = make(map[string][]string, len(paramsList))
	q.Command, _ = SplitKeyValues(paramsList[0])

	for _, param := range paramsList {
		if param == "" {
			continue
		}

		key, values := SplitKeyValues(param)
		q.Params[key] = append(q.Params[key], values...)
	}
}
