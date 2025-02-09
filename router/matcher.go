package router

import (
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
	"github.com/opoccomaxao/tg-instrumentation/texts"
)

type UpdateMatcher func(update *apimodels.Update) bool

type customCommand struct {
	matcher  UpdateMatcher
	handlers []Handler
}

type customCommandList struct {
	commands []*customCommand
}

func (l *customCommandList) AddHandler(
	matcher UpdateMatcher,
	handlers ...Handler,
) {
	l.commands = append(l.commands, &customCommand{
		matcher:  matcher,
		handlers: handlers,
	})
}

func (l *customCommandList) FindHandler(
	update *apimodels.Update,
) ([]Handler, bool) {
	for _, cmd := range l.commands {
		if cmd.matcher(update) {
			return cmd.handlers, true
		}
	}

	return nil, false
}

type command struct {
	matcher  *texts.SimpleMatcher
	pattern  string
	handlers []Handler
}

type commandList struct {
	commands []*command
}

func (l *commandList) AddHandler(
	pattern texts.SimplePattern,
	handlers ...Handler,
) error {
	matcher, err := texts.NewSimpleMatcher(pattern)
	if err != nil {
		//nolint:wrapcheck
		return err
	}

	l.commands = append(l.commands, &command{
		matcher:  matcher,
		pattern:  pattern.String(),
		handlers: handlers,
	})

	return nil
}

func (l *commandList) FindHandler(
	text string,
) ([]Handler, string, bool) {
	var (
		res      *command
		maxScore = -1
	)

	for _, cmd := range l.commands {
		score := cmd.matcher.Match(text)

		if score > maxScore {
			maxScore = score
			res = cmd
		}
	}

	if res == nil {
		return nil, "", false
	}

	return res.handlers, res.pattern, true
}
