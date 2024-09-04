package bot

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

func (b *Bot) HandleMessage(message *echotron.Message) {
	if message.Text == "/start" {
		b.SendStartMessage()
		if !b.SessionSettings.WithDefinitionToWordCards && !b.SessionSettings.WithWordToDefinitionCards {
			b.SessionSettings.WithDefinitionToWordCards = true
			b.SessionSettings.WithWordToDefinitionCards = true
			SaveBotToDb(b.db, b)
		}
	} else if message.Text == reviewCommand {
		b.SendStartSessionMessage()
	} else {
		re := regexp.MustCompile(`^[A-Za-z ]+$`)
		if re.MatchString(message.Text) {
			b.HandleWordMessage(message.Text)
		} else {
			log.Println("Unexpected message got: ", message.Text)
		}
	}
}

func (b *Bot) HandleStartSession(message *echotron.Message) {
	if len(b.StoredUsersWords) == 0 {
		b.SendNoWordsStored()
		return
	}
	b.LoadSessionWords()
	err := b.LoadSessionWordEntries()
	if err != nil {
		log.Println("HandleStartSession error: ", err.Error())
		b.SendWentWrongMessage()
		return
	}
	if len(b.SessionWords) == 0 {
		b.SendNoSessionWords()
		return
	}
	b.sessionWordsQueue = make([]*UsersWord, len(b.SessionWords))
	copy(b.sessionWordsQueue, b.SessionWords)
	b.SendFirstSessionWordMessage()
}

func (b *Bot) HandleWordMessage(word string) {
	wordEntry, err := b.LoadWordEntry(word)

	if err != nil {
		log.Println("HandleWordMessage error: ", err.Error())
		b.SendWentWrongMessage()
		return
	}

	b.SendFirstWordEntryMessage(*wordEntry)
	existsInUsersList := false
	for _, w := range b.StoredUsersWords {
		if w.Word == wordEntry.Word {
			existsInUsersList = true
			break
		}
	}

	if !existsInUsersList {
		usersWord := UsersWord{
			BotId: b.ID,
			Word:  wordEntry.Word,
		}
		usersWord = *insertSessionItemToDbIfNotExists(b.db, usersWord)

		b.StoredUsersWords = append(b.StoredUsersWords, &usersWord)
	}

}

func (b *Bot) HandleSessionFinish(callbackQuery *echotron.CallbackQuery) {
	for _, w := range b.SessionWords {
		if w.reference != nil {
			w.reference.sessionMistakes += w.sessionMistakes
		}
	}
	for _, w := range b.SessionWords {
		w.LastSessionMistakes = w.sessionMistakes
		w.IsNewWord = false
		SaveBotWordEntryInDb(b.db, w)
	}
	b.EditMessageToCompleteSession(callbackQuery.Message)
}

func (b *Bot) HandleCallbackQuery(callbackQuery *echotron.CallbackQuery) {
	b.AnswerCallbackQueryDefault(callbackQuery)
	if callbackQuery.Data == RememberWordButtonData {
		b.sessionWordsQueue[0].rememberRating++
		if b.sessionWordsQueue[0].rememberRating < 1 {
			b.sessionWordsQueue = append(b.sessionWordsQueue, b.sessionWordsQueue[0])
		}
		b.sessionWordsQueue = b.sessionWordsQueue[1:]

		b.NextWord(callbackQuery)

	} else if callbackQuery.Data == NotRememberWordButtonData {
		b.sessionWordsQueue[0].rememberRating--
		if b.sessionWordsQueue[0].rememberRating < -1 {
			b.sessionWordsQueue[0].rememberRating = -1
		}
		b.sessionWordsQueue[0].sessionMistakes++

		b.sessionWordsQueue = append(b.sessionWordsQueue, b.sessionWordsQueue[0])
		b.sessionWordsQueue = b.sessionWordsQueue[1:]

		b.NextWord(callbackQuery)

	} else if callbackQuery.Data == ShowWordButtonData {
		b.EditWordMessageToShowCurrent(callbackQuery.Message)

	} else if callbackQuery.Data == SettingsButtonData {
		b.EditMessageToSetting(callbackQuery)

	} else if callbackQuery.Data == TurnOnNotificationsData {
		b.Notifications = true
		SaveBotToDb(b.db, b)
		b.EditMessageToNotifications(callbackQuery)

	} else if callbackQuery.Data == TurnOffNotificationsData {
		b.Notifications = false
		SaveBotToDb(b.db, b)
		b.EditMessageToNotifications(callbackQuery)

	} else if callbackQuery.Data == GoBackData {
		b.EditLastMessageToStart(callbackQuery)

	} else if callbackQuery.Data == NextWordButtonData {
		b.NextWord(callbackQuery)

	} else if callbackQuery.Data == LearnData {
		b.EditMessageToNewRepeatChoice(callbackQuery.Message)
	} else if callbackQuery.Data == NewWordsButtonData {
		b.SessionSettings.newWords = true
		b.EditMessageToWordsAmount(callbackQuery.Message)

	} else if callbackQuery.Data == RepeatWordsButtonData {
		b.SessionSettings.repeatWords = true
		b.EditMessageToWordsAmount(callbackQuery.Message)

	} else if strings.Contains(callbackQuery.Data, WordsAmountDataLast) {
		re := regexp.MustCompile(`\d+`)
		digits := re.FindAllString(callbackQuery.Data, -1)
		amount, _ := strconv.Atoi(digits[0])
		b.SessionSettings.wordsAmount = amount
		b.HandleStartSession(callbackQuery.Message)

	} else if callbackQuery.Data == NotificationsSettingsData {
		b.EditMessageToNotifications(callbackQuery)
	} else if callbackQuery.Data == SessionSettingsData {
		b.EditMessageToSessionSettings(callbackQuery)

	} else if callbackQuery.Data == WordToDefinitionButtonData {
		b.SessionSettings.WithWordToDefinitionCards = !b.SessionSettings.WithWordToDefinitionCards
		SaveSessionOptionsToDb(b.db, &b.SessionSettings)
		b.EditMessageToSessionSettings(callbackQuery)

	} else if callbackQuery.Data == DefinitionToWordButtonData {
		b.SessionSettings.WithDefinitionToWordCards = !b.SessionSettings.WithDefinitionToWordCards
		SaveSessionOptionsToDb(b.db, &b.SessionSettings)
		b.EditMessageToSessionSettings(callbackQuery)

	} else if callbackQuery.Data == GoBackToSessionSettingsData {
		b.EditMessageToSetting(callbackQuery)

	} else {
		log.Println("Error: Bot got unexpected callback query:", callbackQuery.Data)
		b.SendStartMessage()
	}
}

func (b *Bot) NextWord(callbackQuery *echotron.CallbackQuery) {
	if len(b.sessionWordsQueue) == 0 {
		b.HandleSessionFinish(callbackQuery)
		return
	}
	b.EditWordMessageToCurrent(callbackQuery.Message)
}
