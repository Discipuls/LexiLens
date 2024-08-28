package seeker

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetWordEntry(w http.ResponseWriter, r *http.Request) {
	//word := c.Param("word")
	word := r.URL.Path[len("entry/"):]
	entry := doOnlineGoogleDictionary(word)
	jsonAnswer, err := json.Marshal(entry)
	if err != nil {
		log.Println("Json marshal err: ", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonAnswer)
}
