package bot

import (
	"fmt"
	"strconv"

	"github.com/NicoNex/echotron/v3"
)

func (b *Bot) generateStandartKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         LearnButtonText,
		CallbackData: LearnData,
	})
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         string(settingsButtonText),
		CallbackData: string(SettingsButtonData),
	})

	return keyboard
}

func (b *Bot) generateNotificationsKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], b.generateGoBackButton())
	if b.Notifications {
		keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
			Text:         string(TurnOffNotificationsText),
			CallbackData: TurnOffNotificationsData,
		})
	} else {
		keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
			Text:         string(TurnOnNotificationsText),
			CallbackData: TurnOnNotificationsData,
		})
	}

	return keyboard
}

func (b *Bot) generateGoBackKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], b.generateGoBackButton())

	return keyboard
}

func (b *Bot) generateGoBackButton() echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         string(GoBackButtonText),
		CallbackData: GoBackData,
	}
}

func (b *Bot) generateShowWordKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], b.generateShowWordButton())

	return keyboard
}

func (b *Bot) generateShowWordButton() echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         ShowWordButtonText,
		CallbackData: ShowWordButtonData,
	}
}

func (b *Bot) generateRateWordKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], b.generateNotRememberWordButton())
	keyboard[0] = append(keyboard[0], b.generateRememberWordButton())

	return keyboard
}

func (b *Bot) generateRememberWordButton() echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         RememberWordButtonText,
		CallbackData: RememberWordButtonData,
	}
}
func (b *Bot) generateNotRememberWordButton() echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         NotRememberWordButtonText,
		CallbackData: NotRememberWordButtonData,
	}
}

func (b *Bot) generateNextWordKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], b.generateNextWordButton())

	return keyboard
}

func (b *Bot) generateNextWordButton() echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         NextWordButtonText,
		CallbackData: NextWordButtonData,
	}
}

func (b *Bot) generateStartSessionKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	amounts := []int{
		1, 5, 10, 20, 50,
	}
	for _, amount := range amounts {
		keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
		keyboard[len(keyboard)-1] = append(keyboard[0], b.generateMasterButton(amount))
	}

	return keyboard
}

func (b *Bot) generateMasterButton(count int) echotron.InlineKeyboardButton {
	return echotron.InlineKeyboardButton{
		Text:         fmt.Sprintf(WordsAmountText, strconv.Itoa(count)),
		CallbackData: fmt.Sprintf(WordsAmountData, strconv.Itoa(count)),
	}
}

func (b *Bot) generateLearnKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         LearnButtonText,
		CallbackData: LearnData,
	})
	return keyboard
}

func (b *Bot) GenerateSettingsKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         NotificationsSettingsButtonText,
		CallbackData: NotificationsSettingsData,
	})
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         SessionSettingsButtonText,
		CallbackData: SessionSettingsData,
	})
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[1] = append(keyboard[1], b.generateGoBackButton())
	return keyboard
}

func (b *Bot) GenerateSessionSettingsKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	wordToDefinitionText := WordToDefinitionButtonText
	if b.SessionSettings.WithWordToDefinitionCards {
		wordToDefinitionText += enabledText
	} else {
		wordToDefinitionText += disabledText
	}
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         wordToDefinitionText,
		CallbackData: WordToDefinitionButtonData,
	})

	definitionToWordText := DefinitionToWordButtonText
	if b.SessionSettings.WithDefinitionToWordCards {
		definitionToWordText += enabledText
	} else {
		definitionToWordText += disabledText
	}
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[1] = append(keyboard[1], echotron.InlineKeyboardButton{
		Text:         definitionToWordText,
		CallbackData: DefinitionToWordButtonData,
	})
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[2] = append(keyboard[2], echotron.InlineKeyboardButton{
		Text:         GoBackButtonText,
		CallbackData: GoBackToSessionSettingsData,
	})
	return keyboard
}
