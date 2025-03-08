package router

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
)

// RespondPrivateMessage sends a message to the user from the context.
func (c *Context) RespondPrivateMessage(
	params *bot.SendMessageParams,
) (*models.Message, error) {
	update := c.Update()

	switch {
	case update.Message != nil:
		params.ChatID = update.Message.From.ID
	case update.CallbackQuery != nil:
		params.ChatID = update.CallbackQuery.From.ID
	case update.InlineQuery != nil:
		params.ChatID = update.InlineQuery.From.ID
	default:
		return nil, errors.Wrap(ErrFailed, "unsupported update type")
	}

	return c.SendMessage(params)
}

// RespondReactionEmoji sends a reaction to the message from the context.
func (c *Context) RespondReactionEmoji(
	emoji string,
) (bool, error) {
	update := c.Update()
	if update.Message == nil {
		return false, errors.Wrap(ErrFailed, "supported only messages")
	}

	return c.SetMessageReaction(&bot.SetMessageReactionParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		Reaction: []models.ReactionType{
			{
				Type:              models.ReactionTypeTypeEmoji,
				ReactionTypeEmoji: &models.ReactionTypeEmoji{Emoji: emoji},
			},
		},
	})
}

// RespondPrivateMessageText sends a text message to the user from the context.
func (c *Context) RespondPrivateMessageText(
	text string,
) (*models.Message, error) {
	return c.RespondPrivateMessage(&bot.SendMessageParams{
		Text: text,
	})
}

// RespondCallbackText sends a text response to the callback query from the context.
func (c *Context) RespondCallbackText(
	text string,
) (bool, error) {
	update := c.Update()
	if update.CallbackQuery == nil {
		return false, errors.Wrap(ErrFailed, "supported only callback queries")
	}

	return c.AnswerCallbackQuery(&bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            text,
	})
}

func (c *Context) DeleteMessageFromCallback() (bool, error) {
	update := c.Update()
	if update.CallbackQuery == nil {
		return false, errors.Wrap(ErrFailed, "supported only callback queries")
	}

	params := &bot.DeleteMessageParams{}

	switch {
	case update.CallbackQuery.Message.Message != nil:
		params.ChatID = update.CallbackQuery.Message.Message.Chat.ID
		params.MessageID = update.CallbackQuery.Message.Message.ID
	case update.CallbackQuery.Message.InaccessibleMessage != nil:
		params.ChatID = update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
		params.MessageID = update.CallbackQuery.Message.InaccessibleMessage.MessageID
	default:
		return false, errors.Wrap(ErrFailed, "unsupported message type")
	}

	return c.DeleteMessage(params)
}
