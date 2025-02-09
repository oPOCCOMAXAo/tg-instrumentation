package router

import (
	"github.com/go-telegram/bot"
	"github.com/pkg/errors"
)

func Recover() Handler {
	return func(ctx *Context) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case error:
					ctx.Error(errors.WithStack(v))
				default:
					ctx.Error(errors.Errorf("%+v", v))
				}
			}
		}()

		ctx.Next()
	}
}

func AutoAccept() Handler {
	return func(ctx *Context) {
		defer ctx.Accept()
		ctx.Next()
	}
}

func AutoAnswerCallbackQuery() Handler {
	return func(ctx *Context) {
		defer func() {
			if ctx.IsAccepted() {
				return
			}

			_, err := ctx.AnswerCallbackQuery(&bot.AnswerCallbackQueryParams{
				CallbackQueryID: ctx.Update().CallbackQuery.ID,
			})
			if err != nil {
				ctx.Error(errors.WithStack(err))
			}
		}()

		ctx.Next()
	}
}
