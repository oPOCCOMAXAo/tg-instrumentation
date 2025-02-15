package router

import "github.com/opoccomaxao/tg-instrumentation/query"

// Query returns the query from the text.
// For text message, it returns the query from the text.
// For callback query, it returns the query from the data.
// For inline query, it returns the query from the query.
// For other types of messages, it returns nil.
func (c *Context) Query() *query.Query {
	if c.text == nil {
		return nil
	}

	return query.Decode(*c.text)
}
