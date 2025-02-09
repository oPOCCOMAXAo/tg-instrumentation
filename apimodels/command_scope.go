package apimodels

import "github.com/go-telegram/bot/models"

type CommandScopeType string

const (
	CSDefault               CommandScopeType = "default"
	CSAllPrivateChats       CommandScopeType = "all_private_chats"
	CSAllGroupChats         CommandScopeType = "all_group_chats"
	CSAllChatAdministrators CommandScopeType = "all_chat_administrators"
	CSChat                  CommandScopeType = "chat"
	CSChatAdministrators    CommandScopeType = "chat_administrators"
	CSChatMember            CommandScopeType = "chat_member"
)

func (s CommandScopeType) BotCommandScope() models.BotCommandScope {
	switch s {
	case CSDefault:
		return &models.BotCommandScopeDefault{}
	case CSAllPrivateChats:
		return &models.BotCommandScopeAllPrivateChats{}
	case CSAllGroupChats:
		return &models.BotCommandScopeAllGroupChats{}
	case CSAllChatAdministrators:
		return &models.BotCommandScopeAllChatAdministrators{}
	case CSChat:
		return &models.BotCommandScopeChat{}
	case CSChatAdministrators:
		return &models.BotCommandScopeChatAdministrators{}
	case CSChatMember:
		return &models.BotCommandScopeChatMember{}
	}

	return nil
}
