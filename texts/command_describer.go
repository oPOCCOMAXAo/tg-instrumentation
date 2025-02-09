package texts

import (
	"github.com/opoccomaxao/tg-instrumentation/apimodels"
)

type (
	Scope        = apimodels.CommandScopeType
	LanguageCode = apimodels.LanguageCode
)

type CommandDescription struct {
	Scope        apimodels.CommandScopeType
	LanguageCode LanguageCode
	Description  string
}

type CommandDescriber struct {
	// map[scope]map[language]map[command]description
	data map[Scope]map[LanguageCode]map[string]string
}

func NewCommandDescriber() *CommandDescriber {
	return &CommandDescriber{}
}

func (s *CommandDescriber) AddCommandDescription(
	command string,
	description []CommandDescription,
) {
	for _, value := range description {
		s.addCommandDescriptionSingle(
			command,
			value.Description,
			value.Scope,
			value.LanguageCode,
		)
	}
}

func (s *CommandDescriber) addCommandDescriptionSingle(
	command string,
	description string,
	scope Scope,
	languageCode LanguageCode,
) {
	if description == "" {
		description = command
	}

	if s.data == nil {
		s.data = map[Scope]map[LanguageCode]map[string]string{}
	}

	if s.data[scope] == nil {
		s.data[scope] = map[LanguageCode]map[string]string{}
	}

	if s.data[scope][languageCode] == nil {
		s.data[scope][languageCode] = map[string]string{}
	}

	s.data[scope][languageCode][command] = description
}

func (s *CommandDescriber) ListCommandsParams() []*apimodels.SetMyCommandsParams {
	var res []*apimodels.SetMyCommandsParams

	for scope, next := range s.data {
		for lang, next := range next {
			data := make([]apimodels.BotCommand, 0, len(next))

			for command, description := range next {
				data = append(data, apimodels.BotCommand{
					Command:     command,
					Description: description,
				})
			}

			res = append(res, &apimodels.SetMyCommandsParams{
				Commands:     data,
				Scope:        scope.BotCommandScope(),
				LanguageCode: string(lang),
			})
		}
	}

	return res
}
