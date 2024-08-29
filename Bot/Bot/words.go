package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"gorm.io/gorm"
)

func (b *Bot) LoadSessionWordEntries() error {
	b.sessionWordEntries = make([]*WordEntry, 0)
	for _, w := range b.SessionWords {
		wordEntry, err := b.LoadWordEntry(w.Word)
		if err != nil {
			return errors.New("LoadSessionWordEntries error: " + err.Error())
		}
		b.sessionWordEntries = append(b.sessionWordEntries, wordEntry)
	}
	return nil
}

func (b *Bot) LoadSessionWords(amount uint) {
	b.SessionWords = make([]*UsersWord, 0)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(b.StoredUsersWords), func(i, j int) {
		b.StoredUsersWords[i], b.StoredUsersWords[j] = b.StoredUsersWords[j], b.StoredUsersWords[i]
	})
	for i, w := range b.StoredUsersWords {
		b.StoredUsersWords[i].usedInSession = false
		b.StoredUsersWords[i].sessionMistakes = 0
		b.StoredUsersWords[i].rememberRating = 0
		if w.IsNewWord {
			b.SessionWords = append(b.SessionWords, w)
			b.StoredUsersWords[i].usedInSession = true
			b.StoredUsersWords[i].rememberRating = -1
			amount--
			if amount == 0 {
				break
			}
		}

	}

	sort.Slice(b.StoredUsersWords, func(i int, j int) bool {
		return b.StoredUsersWords[i].LastSessionMistakes >= b.StoredUsersWords[j].LastSessionMistakes
	})
	for i, w := range b.StoredUsersWords {
		if !w.usedInSession {
			if amount == 0 {
				break
			}
			b.SessionWords = append(b.SessionWords, w)
			b.StoredUsersWords[i].usedInSession = true
			amount--
		}

	}
	reversedSessionWords := make([]*UsersWord, 0)
	if b.SessionSettings.WithWordToDefinitionCards {
		reversedSessionWords = make([]*UsersWord, len(b.SessionWords))
		for i, wordEntry := range b.SessionWords {
			reversedEntry := *wordEntry
			reversedEntry.isFrontCard = true
			reversedEntry.reference = wordEntry
			reversedSessionWords[i] = &reversedEntry
		}
	}

	if b.SessionSettings.WithDefinitionToWordCards {
		b.SessionWords = append(b.SessionWords, reversedSessionWords...)
	} else {
		b.SessionWords = reversedSessionWords
	}

	rand.Shuffle(len(b.SessionWords), func(i, j int) {
		b.SessionWords[i], b.SessionWords[j] = b.SessionWords[j], b.SessionWords[i]
	})
	for i := range b.SessionWords {
		b.SessionWords[i].sessionIndex = i
	}
}

func (b *Bot) LoadWordEntry(word string) (*WordEntry, error) {
	wordEntry, dbErr := FindWordEntry(b.db, word)
	if dbErr != nil {
		var err error
		wordEntry, err = b.GetSeekerWordEntry(word)
		if err != nil {
			return &WordEntry{}, errors.New("LoadWordEntry error: " + err.Error())
		}
		if len(wordEntry.SpeechParts) == 0 {
			return &WordEntry{}, errors.New("LoadWordEntry error: seeker empty word got")
		}

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			err = InsertWordEntryToDb(b.db, wordEntry)
			if err != nil {
				log.Println("LoadWordEntry error: cannot insert word entry to db (seeker version returned): ", err.Error())
				return wordEntry, nil
			}
			wordEntryDb, err := FindWordEntry(b.db, word)
			if err != nil {
				return wordEntryDb, errors.New("LoadWordEntry error: cannot load word entry from db  (seeker version returned): " + err.Error())
			}
		} else {
			return wordEntry, errors.New("LoadWordEntry error: cannot load word entry from db  (seeker version returned): " + err.Error())
		}
	}
	return wordEntry, nil
}

func (b *Bot) GetSeekerWordEntry(word string) (*WordEntry, error) {

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
