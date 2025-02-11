package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
	"github.com/opoccomaxao/tg-instrumentation/texts"
)

type Router struct {
	client      *bot.Bot
	debug       bool
	middlewares []Handler
	texts       commandList
	callbacks   commandList
	custom      customCommandList
	describer   *texts.CommandDescriber
	ctxPool     sync.Pool
	bufferPool  sync.Pool
}

func New(opts ...Option) *Router {
	res := &Router{
		describer: texts.NewCommandDescriber(),
	}
	res.ctxPool.New = res.newContext
	res.bufferPool.New = func() any {
		res := &bytes.Buffer{}
		res.Grow(1024) //nolint:mnd

		return res
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (r *Router) Use(handler ...Handler) {
	r.middlewares = append(r.middlewares, handler...)
}

func (r *Router) Text(
	command texts.SimplePattern,
	handler ...Handler,
) TextHandler {
	err := r.texts.AddHandler(command, handler...)
	if err != nil {
		panic(err)
	}

	return &rawHandler{
		pattern:   command.String(),
		describer: r.describer,
	}
}

func (r *Router) Callback(
	command texts.SimplePattern,
	handler ...Handler,
) {
	err := r.callbacks.AddHandler(command, handler...)
	if err != nil {
		panic(err)
	}
}

func (r *Router) Custom(
	matcher UpdateMatcher,
	handler ...Handler,
) {
	r.custom.AddHandler(matcher, handler...)
}

func (r *Router) ListCommandsParams() []*apimodels.SetMyCommandsParams {
	return r.describer.ListCommandsParams()
}

func (r *Router) newContext() any {
	return &Context{
		router: r,
	}
}

func (r *Router) getContext() *Context {
	res := r.ctxPool.Get().(*Context) //nolint:forcetypeassert
	res.reset()

	return res
}

func (r *Router) putContext(ctx *Context) {
	r.ctxPool.Put(ctx)
}

func (r *Router) getBuffer() *bytes.Buffer {
	res := r.bufferPool.Get().(*bytes.Buffer) //nolint:forcetypeassert
	res.Reset()

	return res
}

func (r *Router) putBuffer(b *bytes.Buffer) {
	r.bufferPool.Put(b)
}

func (r *Router) Handle(
	ctx context.Context,
	update *apimodels.Update,
	opts ...ContextOption,
) (bool, error) {
	var (
		ok       bool
		pattern  string
		handlers []Handler
	)

	switch {
	case update.Message != nil:
		handlers, pattern, ok = r.texts.FindHandler(update.Message.Text)
	case update.CallbackQuery != nil:
		handlers, pattern, ok = r.callbacks.FindHandler(update.CallbackQuery.Data)
	}

	if !ok {
		handlers, ok = r.custom.FindHandler(update)
		pattern = "?"
	}

	if !ok {
		return false, nil
	}

	rCtx := r.getContext()
	defer r.putContext(rCtx)

	rCtx.ctx = ctx
	rCtx.update = update
	rCtx.pattern = pattern
	rCtx.handlers = slices.Concat(r.middlewares, handlers)

	for _, opt := range opts {
		opt(rCtx)
	}

	rCtx.Next()

	return rCtx.IsAccepted(), nil
}

// HandlerFunc is Webhook handler as http.HandlerFunc implementation.
func (r *Router) HandlerFunc(
	writer http.ResponseWriter,
	req *http.Request,
) {
	defer req.Body.Close()

	var (
		update apimodels.Update
		opts   []ContextOption
	)

	reader := io.Reader(req.Body)

	if r.debug {
		buf := r.getBuffer()
		defer r.putBuffer(buf)

		reader = io.TeeReader(reader, buf)

		opts = append(opts, func(ctx *Context) {
			ctx.raw = buf
		})
	}

	err := json.NewDecoder(reader).Decode(&update)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	accepted, err := r.Handle(req.Context(), &update, opts...)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	if accepted {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}
