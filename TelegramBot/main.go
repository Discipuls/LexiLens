package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/NicoNex/echotron/v3"
	"gorm.io/gorm"
)

type stateFn func(*echotron.Update) stateFn

type Bot struct {
	ChatID      int64
	state       stateFn
	Words       []WordEntry //TODO declare foreign keys
	CurrentWord int
	gorm.Model

	echotron.API `gorm:"-"`
	db           *gorm.DB
}

type BotCommand string

const (
	startCommand  BotCommand = "/start"
	reviewCommand            = "/r"
)

type BotCallbackData string

const (
	rememberData     BotCallbackData = "next"
	dontRememberData                 = "previous"
)

var token string
var seekerUrl string

type BotController struct {
	db *gorm.DB
}

func (controller BotController) newBot(chatID int64) echotron.Bot {
	defer log.Default().Printf("New bot with chatId: %d\n", chatID)
	b := &Bot{
		chatID,
		nil,
		nil,
		0,
		gorm.Model{},
		echotron.NewAPI(token),
		controller.db,
	}
	b.state = b.handleMessage
	b.Words = make([]WordEntry, 0)
	return b
}

func (b *Bot) Update(update *echotron.Update) {
	b.state = b.state(update)
}

func (b *Bot) handleMessage(update *echotron.Update) stateFn {
	log.Println("handle message")
	if update.Message != nil {
		//log.Println("Starting db connection!")
		//go ConnectDatabase()
		if update.Message.Text == "/start" {
			keyboard := make([][]echotron.InlineKeyboardButton, 0)
			keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
			keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
				Text:         "text",
				CallbackData: "data",
			})
			options := echotron.MessageOptions{
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: keyboard,
				},
			}
			_, err := b.SendMessage("Hello world", b.ChatID, &options)
			if err != nil {
				log.Println("Error: ", err.Error())
			}
		} else if update.Message.Text == reviewCommand {
			return b.handleReview(update)
		} else {
			response, err := http.Get(seekerUrl + "/entry/" + update.Message.Text)
			if err != nil {
				log.Println(err.Error())
				return b.handleMessage
			}
			defer response.Body.Close()
			data, err := io.ReadAll(response.Body)
			if err != nil {
				log.Println(err.Error())
				return b.handleMessage
			}
			var wordEntry WordEntry
			err = json.Unmarshal(data, &wordEntry)
			if err != nil {
				log.Println(err.Error())
				return b.handleMessage
			}
			_, err = b.SendMessage(wordEntry.ToHTML(&EntryFormatOptions{ExamplesLimit: 1, DefinitionsLimit: 2}), update.ChatID(), &echotron.MessageOptions{
				ParseMode: echotron.HTML,
			})
			if err != nil {
				log.Println(err.Error())
				return b.handleMessage
			}
			b.Words = append(b.Words, wordEntry)
		}

	} else if update.CallbackQuery != nil {
		log.Default().Println(update.CallbackQuery.Data)
		b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
	}
	return b.handleMessage
}

func (b *Bot) handleReview(update *echotron.Update) stateFn {

	log.Println("Handle review: ", len(b.Words))
	if update.Message != nil && update.Message.Text == reviewCommand {
		// db, err := ConnectDatabase()
		// if err != nil {
		// 	log.Println(err.Error())
		// } else {
		// 	log.Println("Trying save to db...")
		// 	saveBotToDb(db, b)
		// }

		log.Println("Trying save to db...")
		//saveBotToDbIfNotExists(b.db, b)

		b.CurrentWord = 0
		_, err := b.SendMessage(b.Words[b.CurrentWord].ToHTML(&EntryFormatOptions{ExamplesLimit: 1, DefinitionsLimit: 2}), update.ChatID(),
			&echotron.MessageOptions{
				ParseMode: echotron.HTML,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: CreateReviewKeyboard(),
				},
			})
		if err != nil {
			log.Panic(err.Error())
		}
		return b.handleReview
	} else if update.CallbackQuery != nil {
		defer b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
		if update.CallbackQuery.Data == string(rememberData) || update.CallbackQuery.Data == dontRememberData {
			b.CurrentWord++
			if b.CurrentWord < len(b.Words) {
				_, err := b.SendMessage(b.Words[b.CurrentWord].ToHTML(&EntryFormatOptions{ExamplesLimit: 1, DefinitionsLimit: 2}), update.ChatID(),
					&echotron.MessageOptions{
						ParseMode: echotron.HTML,
						ReplyMarkup: echotron.InlineKeyboardMarkup{
							InlineKeyboard: CreateReviewKeyboard(),
						},
					})
				if err != nil {
					log.Panic(err.Error())
				}
			} else {
				b.CurrentWord = 0
				return b.handleMessage
			}
			return b.handleReview
		} else {
			log.Default().Println(update.CallbackQuery.Data)
			b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
		}
	}

	return b.handleMessage
}

func CreateReviewKeyboard() [][]echotron.InlineKeyboardButton {
	keyboard := make([][]echotron.InlineKeyboardButton, 0)
	keyboard = append(keyboard, make([]echotron.InlineKeyboardButton, 0))
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         "Don't remember",
		CallbackData: dontRememberData,
	})
	keyboard[0] = append(keyboard[0], echotron.InlineKeyboardButton{
		Text:         "Remember",
		CallbackData: string(rememberData),
	})
	return keyboard
}
func main() {
	log.Println("Starting...")
	defer log.Println("Stop")
	configuration, err := GetConfiguration()
	if err != nil {
		panic(err.Error())
	}
	seekerUrl = "http://" + configuration.Seeker.Host + ":" + configuration.Seeker.Port

	token = configuration.Bot.Token

	db, err := ConnectDatabase()
	if err != nil {
		log.Println("Error with database connection: ", err.Error())
	}
	botController := BotController{
		db: db,
	}
	dispatcher := echotron.NewDispatcher(token, botController.newBot)

	for {
		log.Println(dispatcher.Poll())

		time.Sleep(5 * time.Second)
	}
}
