package router

import "github.com/go-telegram/bot"

type Option func(*Router)

func WithClient(client *bot.Bot) Option {
	return func(r *Router) {
		r.client = client
	}
}

func WithDebug() Option {
	return func(r *Router) {
		r.debug = true
	}
}
