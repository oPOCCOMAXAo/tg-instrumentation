package router

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
)

func (c *Context) getClient() (*bot.Bot, error) {
	if c.router.client == nil {
		return nil, errors.Wrap(ErrFailed, "client is not set. use router.New() with router.WithClient() option")
	}

	return c.router.client, nil
}

// AnswerCallbackQuery https://core.telegram.org/bots/api#answercallbackquery
func (c *Context) AnswerCallbackQuery(
	params *bot.AnswerCallbackQueryParams,
) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	res, err := client.AnswerCallbackQuery(c.ctx, params)
	if err != nil {
		return false, errors.WithStack(err)
	}

	c.Accept()

	return res, nil
}

// SendMessage https://core.telegram.org/bots/api#sendmessage
func (c *Context) SendMessage(
	params *bot.SendMessageParams,
) (*models.Message, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	res, err := client.SendMessage(c.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Accept()

	return res, nil
}

// EditMessageText https://core.telegram.org/bots/api#editmessagetext
func (c *Context) EditMessageText(
	params *bot.EditMessageTextParams,
) (*models.Message, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	res, err := client.EditMessageText(c.ctx, params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
