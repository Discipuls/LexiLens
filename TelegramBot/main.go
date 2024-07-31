package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

var token string
var seekerHost string

func newBot(chatID int64) echotron.Bot {
	defer log.Default().Printf("New bot with chatId: %d\n", chatID)
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message != nil {
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
			_, err := b.SendMessage("Hello world", b.chatID, &options)
			if err != nil {
				log.Println("Error: ", err.Error())
			}
		} else {
			respone, err := http.Get(seekerHost + "/entry/" + update.Message.Text)
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer respone.Body.Close()
			data, err := io.ReadAll(respone.Body)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var wordEntry WordEntry
			err = json.Unmarshal(data, &wordEntry)
			if err != nil {
				log.Println(err.Error())
				return
			}
			_, err = b.SendMessage(wordEntry.ToHTML(&EntryFormatOptions{ExamplesLimit: 1, DefinitionsLimit: 2}), update.ChatID(), &echotron.MessageOptions{
				ParseMode: echotron.HTML,
			})
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	} else if update.CallbackQuery != nil {
		log.Default().Println(update.CallbackQuery.Data)
		b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
	}

}

func main() {
	log.Println("Starting...")
	defer log.Println("Stop")
	configuration, err := GetConfiguration()
	if err != nil {
		panic(err.Error())
	}
	seekerHost = ""
	if host := os.Getenv("SEEKER_HOST"); host != "" {
		seekerHost = "http://" + host
	}

	token = configuration.Bot.Token

	dispatcher := echotron.NewDispatcher(token, newBot)

	for {
		log.Println(dispatcher.Poll())

		time.Sleep(5 * time.Second)
	}
}
