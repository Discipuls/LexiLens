package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO multithread
func main() {
	if len(os.Args) > 1 {
		word := os.Args[1]
		entry := doOnlineGoogleDictionary(word)
		printEntry(entry)
	}else{
		router := gin.Default()
		router.GET("/entry/:word", getWordEntry)
		router.Run()
	}
}

func doOnlineGoogleDictionary(word string) *WordEntry {
	url := "https://googledictionary.freecollocation.com/meaning?word=" + word
	body := requestWebPage(url)
	entry, err := ParsedoOnlineGoogleDictionary(body)
	entry.Word = word
	if err != nil {
		panic(err)
	}
	//printEntry(&entry)
	return &entry
}

func printEntry(entry *WordEntry) {
	fmt.Println(entry.Word)
	fmt.Println("-----------------")
	for _, speechPart := range entry.SpeechParts {
		fmt.Println(speechPart.SpeechPart)
		for defIndex, def := range speechPart.Definitions {
			fmt.Println(fmt.Sprintf("\t%d. ", defIndex+1), def.Definition)
			for _, example := range def.Examples {
				fmt.Print("\t\t")
				for _, examplePiece := range example.Pieces {
					if examplePiece.ContainsMainWord {
						fmt.Print("...")
					} else {
						fmt.Print(examplePiece.Value)
					}
				}
				fmt.Println()
			}
		}
	}
}

func requestWebPage(url string) (body []byte) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received non-200 response code: %d\n", response.StatusCode)
	}
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
	}
	return
}
