package router

import (
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
	"github.com/opoccomaxao/tg-instrumentation/texts"
)

type Handler func(ctx *Context)

type TextHandler interface {
	WithDescription(
		language apimodels.LanguageCode,
		scope apimodels.CommandScopeType,
		description string,
	) TextHandler
}

type rawHandler struct {
	pattern   string
	describer *texts.CommandDescriber
}

func (h *rawHandler) WithDescription(
	language apimodels.LanguageCode,
	scope apimodels.CommandScopeType,
	description string,
) TextHandler {
	h.describer.AddCommandDescription(h.pattern, []texts.CommandDescription{
		{
			Scope:        scope,
			LanguageCode: language,
			Description:  description,
		},
	})

	return h
}
