package lexi

import "gorm.io/gorm"

type WordEntry struct {
	gorm.Model
	Word        string            `json:"word"`
	SpeechParts []SpeechPartEntry `json:"speechParts"`
	//BotId       uint
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
