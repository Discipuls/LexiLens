package main

import (
	"fmt"

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
	Definition        []SentencePice     `json:"Definition"`
	Examples          []WordUsageExample `json:"Examples"`
	SpeechPartEntryId uint
}

type WordUsageExample struct {
	gorm.Model
	Pieces           []SentencePice `json:"pieces"`
	WordDefinitionId uint
}

type SentencePice struct {
	gorm.Model
	Value              string `json:"value"`
	ContainsMainWord   bool   `json:"containsMainWord"`
	WordUsageExampleId uint
}

type EntryFormatOptions struct {
	ExamplesLimit    uint
	DefinitionsLimit uint
}

func (entry *WordEntry) ToHTML(options *EntryFormatOptions) string {
	res := ""
	//res += fmt.Sprintf("<tg-spoiler><b><u>%s</u></b></tg-spoiler>\n\n", entry.Word)
	for i, speechPart := range entry.SpeechParts {
		res += fmt.Sprintf("<b><u>%s</u></b>\n\n", speechPart.SpeechPart)

		for defIndex, def := range speechPart.Definitions {
			if defIndex == int(options.DefinitionsLimit) {
				break
			}
			res += fmt.Sprintf("%d) %s\n\n", defIndex+1, def.Definition[0])
			for examplesCount, example := range def.Examples {
				if examplesCount == int(options.ExamplesLimit) {
					break
				}
				res += "    <i>- "
				for _, examplePiece := range example.Pieces {
					if examplePiece.ContainsMainWord {
						res += fmt.Sprintf("<tg-spoiler><b><u>%s</u></b></tg-spoiler>", examplePiece.Value)
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

func (entry *WordEntry) ToString() string {
	res := ""
	res += entry.Word
	res += fmt.Sprintln("-----------------")
	for _, speechPart := range entry.SpeechParts {
		res += fmt.Sprintln(speechPart.SpeechPart)
		for defIndex, def := range speechPart.Definitions {
			res += fmt.Sprintln(fmt.Sprintf("    %d) ", defIndex+1), def.Definition)
			for _, example := range def.Examples {
				res += "        "
				for _, examplePiece := range example.Pieces {
					if examplePiece.ContainsMainWord {
						res += "..."
					} else {
						res += examplePiece.Value
					}
				}
				res += fmt.Sprintln()
			}
		}
	}
	return res
}
