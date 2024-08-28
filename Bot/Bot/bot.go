package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/NicoNex/echotron/v3"
)

var Wg sync.WaitGroup

func (b *Bot) Update(update *echotron.Update) {
	defer Wg.Done()
	if update.CallbackQuery != nil {
		b.HandleCallbackQuery(update.CallbackQuery)
	} else if update.Message != nil {
		b.HandleMessage(update.Message)
	} else {
		log.Println("An unexpected update has been gotten")
	}
}

func (controller BotController) NewBot(chatID int64) echotron.Bot {
	defer log.Default().Printf("New bot with chatId: %d\n", chatID)
	b := insertBotToDbIfNotExists(controller.Db, chatID)
	b.API = echotron.NewAPI(controller.Token)
	b.db = controller.Db
	b.selfUsername = controller.BotUsername
	b.mode = controller.mode
	b.seekerUrl = controller.seekerUrl
	return b
}

func GetSelf(token string) (botAsUser echotron.User, err error) {
	response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token))
	if err != nil {
		return echotron.User{}, errors.New("Error with bot GetMe(): " + err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return echotron.User{}, errors.New("Error with bot GetMe() reading: " + err.Error())
	}

	err = json.Unmarshal([]byte(string(body)[20:len(string(body))-1]), &botAsUser)
	if err != nil {
		return echotron.User{}, errors.New("Error with bot GetMe() unmarshal: " + err.Error())
	}
	return botAsUser, nil
}

func (b *Bot) SendMessage(text string, chatID int64, opts *echotron.MessageOptions) {
	Wg.Add(1)
	go b.sendMessageWg(text, chatID, opts)
}

func (b *Bot) sendMessageWg(text string, chatID int64, opts *echotron.MessageOptions) {
	defer Wg.Done()
	_, err := b.API.SendMessage(text, chatID, opts)
	if err != nil {
		log.Println("Send message error: ", err.Error())
	}

}

func (b *Bot) AnswerCallbackQueryDefault(callbackQuery *echotron.CallbackQuery) {
	Wg.Add(1)
	go b.answerCallbackDefaultWg(callbackQuery)
}

func (b *Bot) answerCallbackDefaultWg(callbackQuery *echotron.CallbackQuery) {
	defer Wg.Done()
	_, err := b.AnswerCallbackQuery(callbackQuery.ID, nil)
	if err != nil {
		log.Println("Answer callback default error:", err.Error())
	}
}

func (b *Bot) DeleteMessage(chatID int64, messageID int) {
	Wg.Add(1)
	go b.deleteMessageWg(chatID, messageID)
}

func (b *Bot) deleteMessageWg(chatID int64, messageID int) {
	defer Wg.Done()
	_, err := b.API.DeleteMessage(chatID, messageID)
	if err != nil {
		log.Println("DeleteMessage error: ", err.Error())
	}
}

func (b *Bot) EditMessageText(text string, msg echotron.MessageIDOptions, opts *echotron.MessageTextOptions) {
	Wg.Add(1)
	go b.editMessageTextWg(text, msg, opts)
}

func (b *Bot) editMessageTextWg(text string, msg echotron.MessageIDOptions, opts *echotron.MessageTextOptions) {
	defer Wg.Done()
	_, err := b.API.EditMessageText(text, msg, opts)
	if err != nil {
		log.Println("EditMessageText error: ", err.Error())
	}
}

func (b *Bot) EditMessageReplyMarkup(msg echotron.MessageIDOptions, opts *echotron.MessageReplyMarkupOptions) {
	Wg.Add(1)
	go b.editMessageReplyMarkupWg(msg, opts)
}

func (b *Bot) editMessageReplyMarkupWg(msg echotron.MessageIDOptions, opts *echotron.MessageReplyMarkupOptions) {
	defer Wg.Done()
	_, err := b.API.EditMessageReplyMarkup(msg, opts)
	if err != nil {
		log.Println("EditMessageReplyMarkup error: ", err.Error())
	}
}

func (b *Bot) editMessageOrSendWg(text string, messageId int, opts *echotron.MessageOptions) {
	defer Wg.Done()
	inlineKeyboardMarkup, _ := opts.ReplyMarkup.(echotron.InlineKeyboardMarkup)
	_, err := b.API.EditMessageText(text, echotron.NewMessageID(b.ChatID, messageId), &echotron.MessageTextOptions{
		ParseMode:   opts.ParseMode,
		ReplyMarkup: inlineKeyboardMarkup,
	})
	if err == nil || err.Error() == `API error: 400 Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message` {
		return
	} else {
		log.Printf("Error: trying to edit message: error: \"%s\"", err.Error())
		b.SendMessage(text, b.ChatID, opts)
		return
	}
}

func (b *Bot) EditLastMessageOrSend(text string, messageId int, opts *echotron.MessageOptions) {
	Wg.Add(1)
	go b.editMessageOrSendWg(text, messageId, opts)
}
