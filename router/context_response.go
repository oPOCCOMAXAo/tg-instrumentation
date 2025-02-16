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
