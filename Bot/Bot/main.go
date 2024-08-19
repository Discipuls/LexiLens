package bot

import (
	"errors"
	"log"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func StartBot() (dispatcher *echotron.Dispatcher, err error) {
	log.Println("greenton-telegram-user-bot is starting...")

	configuration, err := GetConfiguration()
	if err != nil {
		return &echotron.Dispatcher{}, errors.New("Error getting configuration: " + err.Error())
	}

	database, err := ConnectDatabase(&configuration)
	if err != nil {
		return &echotron.Dispatcher{}, errors.New("Error getting database: " + err.Error())
	}
	botAsUser, err := GetSelf(configuration.Bot.Token)

	if err != nil {
		return &echotron.Dispatcher{}, errors.New("Error getting bot self: " + err.Error())
	}

	botController := BotController{
		Token:       configuration.Bot.Token,
		MiniAppURL:  configuration.MiniApp.Url,
		Db:          database,
		BotUsername: botAsUser.Username,
		mode:        configuration.Mode,
		seekerUrl:   configuration.Seeker.Host + configuration.Seeker.Port,
	}

	dispatcher = echotron.NewDispatcher(botController.Token, botController.NewBot)
	return dispatcher, nil
}

func Main() {
	dispatcher, err := StartBot()
	if err != nil {
		log.Panicln("Error starting bot: ", err.Error())
	} else {
		for {
			log.Println(dispatcher.Poll())

			time.Sleep(5 * time.Second)
		}
	}

}
