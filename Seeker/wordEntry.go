package main

import (
	"errors"
	"strings"
)

type WordEntry struct {
	Word        string            `json:"word"`
	SpeechParts []SpeechPartEntry `json:"speechParts"`
}

func (wordEntry *WordEntry) addSpeechPart(speechPart *SpeechPartEntry) error {
	wordEntry.SpeechParts = append(wordEntry.SpeechParts, *speechPart)
	return nil
}

func (wordEntry *WordEntry) addDefinition(definition *WordDefinition) error {
	if len(wordEntry.SpeechParts) == 0 {
		return errors.New("WordEntry error: attempting to add definition to entry with 0 speechParts")
	}
	wordDefinitions := &wordEntry.SpeechParts[len(wordEntry.SpeechParts)-1].Definitions

	*wordDefinitions = append(*wordDefinitions, *definition)
	return nil
}

func (wordEntry *WordEntry) addWordUsageExample() error {
	if len(wordEntry.SpeechParts) == 0 {
		return errors.New("WordEntry error: attempting to add wordUsageExample to entry with 0 speechParts")
	}
	wordDefinitions := &wordEntry.SpeechParts[len(wordEntry.SpeechParts)-1].Definitions
	if len(*wordDefinitions) == 0 {
		return errors.New("WordEntry error: attempting to add wordUsageExample to SpeechPart with 0 wordDefinitions")
	}
	examples := &(*wordDefinitions)[len(*wordDefinitions)-1].Examples
	*examples = append(*examples, WordUsageExample{})
	return nil
}

func (wordEntry *WordEntry) addWordExamplePieces(pieces []WordExamplePice) error {
	if len(wordEntry.SpeechParts) == 0 {
		return errors.New("WordEntry error: attempting to add wordExamplePieces to entry with 0 speechParts")
	}
	wordDefinitions := &wordEntry.SpeechParts[len(wordEntry.SpeechParts)-1].Definitions
	if len(*wordDefinitions) == 0 {
		return errors.New("WordEntry error: attempting to add wordExamplePieces to SpeechPart with 0 wordDefinitions")
	}
	examples := &(*wordDefinitions)[len(*wordDefinitions)-1].Examples
	if len(*examples) == 0 {
		return errors.New("WordEntry error: attempting to add wordExamplePieces to wordDefinitions with 0 wordExamples")
	}

	normalizeWordExamplePieces(&pieces)
	wordExamplePieces := &(*examples)[len(*examples)-1].Pieces
	for _, piece := range pieces {
		adding := false
		for _, el := range piece.Value {
			if el != ' ' {
				adding = true
				break
			}
		}
		if adding {
			*wordExamplePieces = append(*wordExamplePieces, piece)
		}
	}
	return nil
}

func normalizeWordExamplePieces(pieces* []WordExamplePice){
	FirstExamplePiece := &(*pieces)[0]
	if FirstExamplePiece.Value[0] == '-' {
		FirstExamplePiece.Value = strings.Replace(FirstExamplePiece.Value, "-", "", 1)
	}
	if FirstExamplePiece.Value[0] == ' ' {
		FirstExamplePiece.Value = strings.Replace(FirstExamplePiece.Value, " ", "", 1)
	}

	LastExamplePiece := &(*pieces)[len(*pieces)-1]
	if len(LastExamplePiece.Value) > 1 && LastExamplePiece.Value[len(LastExamplePiece.Value)-1] == ' ' {
		LastExamplePiece.Value = LastExamplePiece.Value[:len(LastExamplePiece.Value)-1]
	}
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
