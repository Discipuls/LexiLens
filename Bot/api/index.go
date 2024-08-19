package handler

import (
	"log"
	"net/http"

	bot "github.com/Discipuls/LexiLens/Bot/Bot"
	"github.com/NicoNex/echotron/v3"
)

var dsp *echotron.Dispatcher

func Handler(w http.ResponseWriter, r *http.Request) {
	token, err := bot.GetToken()
	if err != nil {
		log.Panicln("Couldn't get token: ", err.Error())
	}

	if r.URL.Path == ("/" + token) {
		if dsp == nil {
			dispatcher, err := bot.StartBot()
			if err != nil {
				log.Panicln("Error starting bot: ", err.Error())
			} else {
				dsp = dispatcher
			}
		}
		log.Println("handle")
		bot.Wg.Add(1)
		dsp.HandleWebhook(w, r)
		bot.Wg.Wait()
	} else {
		w.Write([]byte("Wrong url"))
	}
}
