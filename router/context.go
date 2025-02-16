package router

import (
	"bytes"
	"context"
	"math"

	"github.com/go-telegram/bot/models"
)

const (
	abortIndex = math.MaxInt16 >> 1
)

type ContextOption func(*Context)

type Context struct {
	CtxTemp

	router *Router
}

type CtxTemp struct {
	ctx      context.Context //nolint:containedctx
	update   *models.Update
	text     *string
	pattern  string
	raw      *bytes.Buffer
	handlers []Handler
	errors   []error
	index    int
	accepted bool
}

func (c *Context) reset() {
	c.CtxTemp = CtxTemp{
		index: -1,
	}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Update() *models.Update {
	return c.update
}

func (c *Context) Pattern() string {
	return c.pattern
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++

	for c.index < len(c.handlers) {
		h := c.handlers[c.index]
		if h != nil {
			h(c)
		}

		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails, call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Accept marks update as accepted.
// If the update is accepted, the router will not try to process it again.
func (c *Context) Accept() {
	c.accepted = true
}

func (c *Context) IsAccepted() bool {
	return c.accepted
}

func (c *Context) Error(err error) {
	c.errors = append(c.errors, err)
}

func (c *Context) Errors() []error {
	return c.errors
}

func (c *Context) RawDebug() []byte {
	if c.raw == nil {
		return nil
	}

	return c.raw.Bytes()
}

// LogError1 checks and logs an error.
//
// Use it if you don't care about the error, but want to log it.
//
// Example:
//
//	func doSomething() error {...}
//
//	ctx.LogError1(doSomething())
func (c *Context) LogError1(err error) {
	if err != nil {
		c.Error(err)
	}
}

// LogError2 checks and logs an error.
//
// Use it if you don't care about the error, but want to log it.
//
// Example:
//
//	func doSomething2() (sometype, error) {...}
//
//	ctx.LogError2(doSomething2())
func (c *Context) LogError2(_ any, err error) {
	if err != nil {
		c.Error(err)
	}
}

// LogError3 checks and logs an error.
//
// Use it if you don't care about the error, but want to log it.
//
// Example:
//
//	func doSomething3() (sometype, sometype2, error) {...}
//
//	ctx.LogError3(doSomething3())
func (c *Context) LogError3(_ any, _ any, err error) {
	if err != nil {
		c.Error(err)
	}
}
