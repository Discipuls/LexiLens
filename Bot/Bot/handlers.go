package bot

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
	"gorm.io/gorm"
)

func (b *Bot) HandleMessage(message *echotron.Message) {
	if message.Text == "/start" {
		b.SendStartMessage()
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

func (b *Bot) HandleStartSession(message *echotron.Message, amount uint) {
	b.currentWord = 0
	b.LoadSessionWords(amount)
	b.sessionsWordCopy = make([]*BotWordEntry, len(b.sessionWords))
	copy(b.sessionsWordCopy, b.sessionWords)
	b.EditWordMessageToCurrent(message)
}

func (b *Bot) HandleWordMessage(word string) {
	wordEntry, dbErr := FindWordEntry(b.db, word)
	if dbErr != nil {
		var err error
		wordEntry, err = b.GetWordEntry(word)
		if err != nil {
			log.Println("HandleWordMessage error:", err.Error())
			b.SendWentWrongMessage()

			return
		}
		if len(wordEntry.SpeechParts) == 0 {
			log.Println("Word not found")
			b.SendWentWrongMessage()
		}

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			err = InsertWordEntryToDb(b.db, wordEntry)
			if err != nil {
				log.Println("HandleWordMessage error:", err.Error())
				return
			}
			wordEntry, err = FindWordEntry(b.db, word)
			if err != nil {
				log.Println("HandleWordMessage error:", err.Error())

			}
		}
	}
	b.SendFirstWordEntryMessage(*wordEntry)
	existsInUsersList := false
	for _, w := range b.WordEntries {
		if w.Word == wordEntry.Word {
			existsInUsersList = true
			break
		}
	}

	if !existsInUsersList {
		botWordEntry := BotWordEntry{
			BotId: b.ID,
			Word:  wordEntry.Word,
		}
		botWordEntry = *insertBotWordEntryToDbIfNotExists(b.db, botWordEntry)

		b.WordEntries = append(b.WordEntries, &botWordEntry)
	}

}

func (b *Bot) HandleSessionFinish(callbackQuery *echotron.CallbackQuery) {
	for _, w := range b.sessionsWordCopy {
		w.LastSessionMistakes = w.sessionMistakes
		w.IsNewWord = false
		SaveBotWordEntryInDb(b.db, w)
	}
	b.EditMessageToCompleteSession(callbackQuery.Message)
	//b.DeleteMessage(b.ChatID, callbackQuery.Message.ID)
}

func (b *Bot) HandleCallbackQuery(callbackQuery *echotron.CallbackQuery) {
	b.AnswerCallbackQueryDefault(callbackQuery)
	if callbackQuery.Data == RememberWordButtonData {
		b.sessionWords[b.currentWord].rememberRating++
		if b.sessionWords[b.currentWord].rememberRating > 0 {
			b.sessionWords = removeBotWordEntry(b.sessionWords, int(b.currentWord))
			b.currentWord--
		}

		b.NextWord(callbackQuery)

	} else if callbackQuery.Data == NotRememberWordButtonData {
		b.sessionWords[b.currentWord].rememberRating--
		b.sessionWords[b.currentWord].sessionMistakes++
		b.NextWord(callbackQuery)

	} else if callbackQuery.Data == ShowWordButtonData {
		//	log.Println(b.sessionWords[b.currentWord].rememberRating)
		b.EditWordMessageToShowCurrent(callbackQuery.Message)

	} else if callbackQuery.Data == SettingsData {
		b.EditLastToNotifications(callbackQuery)

	} else if callbackQuery.Data == TurnOnNotificationsData {
		b.Notifications = true
		SaveBotToDb(b.db, b)
		b.EditLastToNotifications(callbackQuery)

	} else if callbackQuery.Data == TurnOffNotificationsData {
		b.Notifications = false
		SaveBotToDb(b.db, b)
		b.EditLastToNotifications(callbackQuery)

	} else if callbackQuery.Data == GoBackData {
		b.EditLastMessageToStart(callbackQuery)

	} else if callbackQuery.Data == NextWordButtonData {
		b.NextWord(callbackQuery)
	} else if callbackQuery.Data == LearnData {
		b.EditMessageToStartSession(callbackQuery.Message)
	} else if strings.Contains(callbackQuery.Data, WordsAmountDataLast) {
		re := regexp.MustCompile(`\d+`)
		digits := re.FindAllString(callbackQuery.Data, -1)
		amount, _ := strconv.Atoi(digits[0])
		log.Println(amount)
		b.HandleStartSession(callbackQuery.Message, uint(amount))
	} else {
		log.Println("Error: Bot got unexpected callback query:", callbackQuery.Data)
		b.SendStartMessage()
	}
}

func (b *Bot) NextWord(callbackQuery *echotron.CallbackQuery) {
	b.currentWord++
	if len(b.sessionWords) == 0 {
		b.HandleSessionFinish(callbackQuery)
		return
	} else if int(b.currentWord) == len(b.sessionWords) {
		b.currentWord = 0
	}
	log.Println(b.currentWord)
	b.EditWordMessageToCurrent(callbackQuery.Message)
}
