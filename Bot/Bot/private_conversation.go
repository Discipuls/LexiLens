package bot

import (
	"log"

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

func (b *Bot) EditMessageToNotifications(callbackQuery *echotron.CallbackQuery) {
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

func (b *Bot) EditMessageToSetting(callbackQuery *echotron.CallbackQuery) {
	keyboardMarkup := echotron.InlineKeyboardMarkup{
		InlineKeyboard: b.GenerateSettingsKeyboard(),
	}
	b.EditLastMessageOrSend(settingsText, callbackQuery.Message.ID, &echotron.MessageOptions{
		ReplyMarkup: keyboardMarkup,
	})
}

func (b *Bot) EditMessageToSessionSettings(callbackQuery *echotron.CallbackQuery) {
	keyboardMarkup := echotron.InlineKeyboardMarkup{
		InlineKeyboard: b.GenerateSessionSettingsKeyboard(),
	}
	b.EditLastMessageOrSend(sessionSettingsText, callbackQuery.Message.ID, &echotron.MessageOptions{
		ReplyMarkup: keyboardMarkup,
	})
}

func (b *Bot) SendFirstWordEntryMessage(wordEntry WordEntry) {
	messageText := wordEntry.ToHTML(&EntryFormatOptions{
		ExamplesLimit:    2,
		DefinitionsLimit: 2,
		IsWordHidden:     false,
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
	if int(b.currentWordIndex) >= len(b.sessionWords) {
		log.Println("Error: Current word index is more then session words amount")
		b.SendWentWrongMessage() //TODO fix
	}
	wordEntry := b.sessionWordEntries[b.currentWordIndex]

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

func (b *Bot) SendFirstSessionWordMessage() {
	if len(b.sessionWords) == 0 {
		log.Println("SendFirstSessionWordMessage error: b.sessionWords is empty")
		return
	}
	currentBotWordEntry := b.sessionWords[b.currentWordIndex]

	if currentBotWordEntry.isFrontCard {
		b.SendMessage(currentBotWordEntry.Word, b.ChatID,
			&echotron.MessageOptions{
				ParseMode: echotron.HTML,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: b.generateShowWordKeyboard(),
				},
			})
	} else {
		wordEntry := b.sessionWordEntries[0]

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
}

func (b *Bot) EditWordMessageToCurrent(message *echotron.Message) {
	if len(b.sessionWords) == 0 {
		log.Println("EditWordMessageToCurrent error: b.sessionWords is empty")
		return
	}
	currentBotWordEntry := b.sessionWords[b.currentWordIndex]
	if currentBotWordEntry.isFrontCard {
		b.EditMessageText(currentBotWordEntry.Word, echotron.NewMessageID(b.ChatID, message.ID),
			&echotron.MessageTextOptions{
				ParseMode: echotron.HTML,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: b.generateShowWordKeyboard(),
				},
			})
	} else {
		wordEntry := b.sessionWordEntries[b.currentWordIndex]

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

func (b *Bot) SendWentWrongMessage() {
	b.SendMessage(somethingWentWrongText, b.ChatID, nil)
}

func (b *Bot) SendNoWordsStored() {
	b.SendMessage(noStoredWordsText,
		b.ChatID, nil)
}

func (b *Bot) SendNoSessionWords() {
	b.SendMessage(noWordsWithFiltersText,
		b.ChatID, nil)
}

func (b *Bot) EditMessageToCompleteSession(message *echotron.Message) {
	b.EditMessageText(sessionCompletedText,
		echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStandartKeyboard(),
			},
		})
}

func (b *Bot) SendStartSessionMessage() {
	b.SendMessage(howManyWordsText,
		b.ChatID,
		&echotron.MessageOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStartSessionKeyboard(),
			},
		})
}
func (b *Bot) EditMessageToStartSession(message *echotron.Message) {
	b.EditMessageText(howManyWordsText,
		echotron.NewMessageID(b.ChatID, message.ID),
		&echotron.MessageTextOptions{
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: b.generateStartSessionKeyboard(),
			},
		})
}
