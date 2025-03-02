package router

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
)

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
