package handler

import (
	"net/http"
	"strings"

	seeker "github.com/Discipuls/LexiLensCLI/Seeker"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/word/") {
		seeker.GetWordEntry(w, r)
	} else {
		w.Write([]byte("Wrong request path"))
	}
}
