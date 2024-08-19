package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"gorm.io/gorm"
)

type WordEntry struct {
	gorm.Model
	Word        string            `json:"word"`
	SpeechParts []SpeechPartEntry `json:"speechParts"`
	BotId       uint
}

type SpeechPartEntry struct {
	gorm.Model
	SpeechPart  string           `json:"SpeechPart"`
	Definitions []WordDefinition `json:"Definitions"`
	WordEntryId uint
}

type WordDefinition struct {
	gorm.Model
	Definition        []DefinitionPiece  `json:"Definition" gorm:"foreignKey:DefinitionId;references:ID"`
	WordUsageExamples []WordUsageExample `json:"Examples" gorm:"foreignKey:WordDefinitionId;references:ID"`
	SpeechPartEntryId uint
}

type WordUsageExample struct {
	gorm.Model
	Pieces           []SentencePice `json:"pieces"  gorm:"foreignKey:WordUsageExampleId;references:ID"`
	WordDefinitionId uint
}

type SentencePice struct {
	gorm.Model
	Value              string `json:"value"`
	ContainsMainWord   bool   `json:"containsMainWord"`
	WordUsageExampleId uint
}

type DefinitionPiece struct {
	gorm.Model
	Value            string `json:"value"`
	ContainsMainWord bool   `json:"containsMainWord"`
	DefinitionId     uint
}

type EntryFormatOptions struct {
	ExamplesLimit    uint
	DefinitionsLimit uint
	IsWordHidden     bool
}

type BotWordEntry struct {
	gorm.Model
	BotId               uint
	Word                string
	rememberRating      int
	sessionMistakes     uint
	usedInSession       bool
	LastSessionMistakes uint
	IsNewWord           bool `gorm:"default:true"`
}

func removeBotWordEntry(slice []*BotWordEntry, index int) []*BotWordEntry {
	if index == 0 {
		return slice[index+1:]
	}
	return append(slice[:index], slice[index+1:]...)
}

func (b *Bot) LoadSessionWords(amount uint) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(b.WordEntries), func(i, j int) {
		b.WordEntries[i], b.WordEntries[j] = b.WordEntries[j], b.WordEntries[i]
	})

	for i, w := range b.WordEntries {
		b.WordEntries[i].usedInSession = false
		b.WordEntries[i].sessionMistakes = 0
		b.WordEntries[i].rememberRating = 0
		if w.IsNewWord {
			b.sessionWords = append(b.sessionWords, w)
			b.WordEntries[i].usedInSession = true
			b.WordEntries[i].rememberRating = -1
			amount--
			if amount == 0 {
				return
			}
		}

	}

	sort.Slice(b.WordEntries, func(i int, j int) bool {
		return b.WordEntries[i].LastSessionMistakes > b.WordEntries[j].LastSessionMistakes
	})
	for i, w := range b.WordEntries {
		if !w.usedInSession {
			b.sessionWords = append(b.sessionWords, w)
			b.WordEntries[i].usedInSession = true
			amount--
			if amount == 0 {
				return
			}
		}

	}
}

func (b *Bot) GetWordEntry(word string) (*WordEntry, error) {

	retryCount := 3
	var err error
	var wordEntry WordEntry

	for i := 0; i < retryCount; i += 1 {
		var response *http.Response
		response, err = http.Get(b.seekerUrl + "/word/" + word)
		if err != nil {
			continue
		}
		defer response.Body.Close()
		var data []byte
		data, err = io.ReadAll(response.Body)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		err = json.Unmarshal(data, &wordEntry)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	if err != nil {
		return &WordEntry{}, err
	}

	return &wordEntry, nil
}

func (entry *WordEntry) ToHTML(options *EntryFormatOptions) string {
	res := ""
	if !options.IsWordHidden {
		res += fmt.Sprintf("<b><u>%s</u></b>\n\n", entry.Word)
	}
	for i, speechPart := range entry.SpeechParts {
		res += fmt.Sprintf("<b><u>%s</u></b>\n\n", speechPart.SpeechPart)

		for defIndex, def := range speechPart.Definitions {
			if defIndex == int(options.DefinitionsLimit) {
				break
			}
			definitionText := ""
			for _, definitionPiece := range def.Definition {
				if !definitionPiece.ContainsMainWord {
					definitionText += definitionPiece.Value
				} else {
					if options.IsWordHidden {
						definitionText += "░░░"
					} else {
						definitionText += fmt.Sprintf("<b><u>%s</u></b>", definitionPiece.Value)
					}
				}
			}
			res += fmt.Sprintf("%d) %s\n\n", defIndex+1, definitionText)
			for examplesCount, example := range def.WordUsageExamples {
				if examplesCount == int(options.ExamplesLimit) {
					break
				}
				res += "    <i>- "
				for _, examplePiece := range example.Pieces {
					if examplePiece.ContainsMainWord {
						if options.IsWordHidden {
							res += "</i>░░░ <i>"
						} else {
							res += fmt.Sprintf("<b><u>%s</u></b>", examplePiece.Value)
						}
					} else {
						res += examplePiece.Value
					}
				}
				res += "</i>\n\n"
			}
		}
		if i != len(entry.SpeechParts)-1 {
			res += "<s>__________________________________________</s>\n\n"
		}
	}
	return res
}
