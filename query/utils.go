package query

import "strings"

func JoinKeyValues(key string, values []string) string {
	if len(values) == 0 {
		return key
	}

	return key + QueryValueDelimiter + strings.Join(values, QuerySliceDelimiter)
}

//nolint:nonamedreturns,mnd
func SplitKeyValues(keyValues string) (key string, values []string) {
	parts := strings.SplitN(keyValues, QueryValueDelimiter, 2)
	if len(parts) >= 1 {
		key = parts[0]
	}

	if len(parts) >= 2 {
		values = strings.Split(parts[1], QuerySliceDelimiter)
	}

	return key, values
}
