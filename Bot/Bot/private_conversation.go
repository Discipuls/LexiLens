package bot

import (
	"fmt"
	"log"
	"strconv"

	"github.com/NicoNex/echotron/v3"
)

func (b *Bot) SendStartMessage() {
	b.SendMessage(startMessage, b.ChatID, &echotron.MessageOptions{
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: b.generateStandartKeyboard(),
		},
	})
}

func (b *Bot) EditLastMessageToStart(callbackQuery *echotron.CallbackQuery) {
	b.EditLastMessageOrSend(string(startMessage), callbackQuery.Message.ID, &echotron.MessageOptions{
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: b.generateStandartKeyboard(),
		},
	})
}

func (b *Bot) EditLastToNotifications(callbackQuery *echotron.CallbackQuery) {
	keyboardMarkup := echotron.InlineKeyboardMarkup{
		InlineKeyboard: b.generateNotificationsKeyboard(),
	}
	notificationMessage := ""
	if b.Notifications {
		notificationMessage = notificationsEnabledMessage
	} else {
		notificationMessage = notificationsDisabledMessage
	}

	b.EditLastMessageOrSend(string(notificationMessage), callbackQuery.Message.ID, &echotron.MessageOptions{
		ParseMode:   echotron.HTML,
		ReplyMarkup: keyboardMarkup,
	})
}

func (b *Bot) generateStandartKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         LearnButtonText,
		CallbackData: LearnData, //TODO move to consts
	})
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         string(settingsButtonText),
		CallbackData: string(SettingsData),
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

func (b *Bot) SendFirstWordEntryMessage(wordEntry WordEntry) {
	messageText := wordEntry.ToHTML(&EntryFormatOptions{
		ExamplesLimit:    2,
		DefinitionsLimit: 2,
		IsWordHidden:     true,
	})

	b.SendMessage(messageText, b.ChatID,
		&echotron.MessageOptions{
			ParseMode: echotron.HTML,
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateLearnKeyboard(),
			},
		})
}

func (b *Bot) SendReviewWordEntryMessage(wordEntry WordEntry) {
	messageText := wordEntry.ToHTML(&EntryFormatOptions{
		ExamplesLimit:    2,
		DefinitionsLimit: 2,
		IsWordHidden:     true,
	})

	b.SendMessage(messageText, b.ChatID,
		&echotron.MessageOptions{
			ParseMode: echotron.HTML,
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateShowWordKeyboard(),
			},
		})
}

func (b *Bot) EditWordMessageToShowCurrent(message *echotron.Message) {
	if int(b.currentWord) >= len(b.WordEntries) {
		b.SendWentWrongMessage() //TODO fix
	}
	wordEntry, err := FindWordEntry(b.db, b.sessionWords[b.currentWord].Word)
	if err != nil {
		log.Println("EditWordMessageToNext error: ", err.Error())
		return
	}
	messageText := wordEntry.ToHTML(&EntryFormatOptions{
		ExamplesLimit:    2,
		DefinitionsLimit: 2,
		IsWordHidden:     false,
	})

	b.EditMessageText(messageText, echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ParseMode: echotron.HTML,
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateRateWordKeyboard(),
			},
		})
}

func (b *Bot) EditWordMessageToCurrent(message *echotron.Message) {
	wordEntry, err := FindWordEntry(b.db, b.sessionWords[b.currentWord].Word)
	if err != nil {
		log.Println("EditWordMessageToNext error: ", err.Error())
		return
	}
	messageText := wordEntry.ToHTML(&EntryFormatOptions{
		ExamplesLimit:    2,
		DefinitionsLimit: 2,
		IsWordHidden:     true,
	})

	b.EditMessageText(messageText, echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ParseMode: echotron.HTML,
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateShowWordKeyboard(),
			},
		})
}

func (b *Bot) EditWordReviewKeyboardToNext(message *echotron.Message) {
	b.EditMessageReplyMarkup(echotron.NewMessageID(b.ChatID, message.ID), &echotron.MessageReplyMarkupOptions{
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: b.generateNextWordKeyboard(),
		},
	})
}

func (b *Bot) EditWordReviewKeyboardToShow(message *echotron.Message) {
	b.EditMessageReplyMarkup(echotron.NewMessageID(b.ChatID, message.ID), &echotron.MessageReplyMarkupOptions{
		ReplyMarkup: echotron.InlineKeyboardMarkup{
			InlineKeyboard: b.generateShowWordKeyboard(),
		},
	})
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

func (b *Bot) SendWentWrongMessage() {
	b.SendMessage("Something went wrong(", b.ChatID, nil)
}

func (b *Bot) EditMessageToCompleteSession(message *echotron.Message) {
	b.EditMessageText("Congratulations🎉\n\nYou've completed your study session",
		echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStandartKeyboard(),
			},
		})
}

func (b *Bot) SendStartSessionMessage() {
	b.SendMessage("How many words would you like to master?",
		b.ChatID,
		&echotron.MessageOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStartSessionKeyboard(),
			},
		})
}

func (b *Bot) EditMessageToStartSession(message *echotron.Message) {
	b.EditMessageText("How many words would you like to master?",
		echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStartSessionKeyboard(),
			},
		})
}

func (b *Bot) generateStartSessionKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	amounts := []int{
		0, 1, 2, 50,
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
