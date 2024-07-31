package main

import "fmt"

type WordEntry struct {
	Word        string            `json:"word"`
	SpeechParts []SpeechPartEntry `json:"speechParts"`
}

type SpeechPartEntry struct {
	SpeechPart  string           `json:"SpeechPart"`
	Definitions []WordDefinition `json:"Definitions"`
}

type WordDefinition struct {
	Definition string             `json:"Definition"`
	Examples   []WordUsageExample `json:"Examples"`
}

type WordUsageExample struct {
	Pieces []WordExamplePice `json:"pieces"`
}

type WordExamplePice struct {
	Value            string `json:"value"`
	ContainsMainWord bool   `json:"containsMainWord"`
}

type EntryFormatOptions struct {
	ExamplesLimit    uint
	DefinitionsLimit uint
}

func (entry *WordEntry) ToHTML(options *EntryFormatOptions) string {
	res := ""
	//res += fmt.Sprintf("<tg-spoiler><b><u>%s</u></b></tg-spoiler>\n\n", entry.Word)
	for _, speechPart := range entry.SpeechParts {
		res += fmt.Sprintf("<b>ㅤㅤㅤㅤㅤ%s</b>\n", speechPart.SpeechPart)

		for defIndex, def := range speechPart.Definitions {
			if defIndex == int(options.DefinitionsLimit) {
				break
			}
			res += fmt.Sprintf("%d) %s\n\n", defIndex+1, def.Definition)
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
		res += "\n"
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
