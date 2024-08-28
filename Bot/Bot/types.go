package bot

import (
	"github.com/NicoNex/echotron/v3"
	"gorm.io/gorm"
)

type BotController struct {
	Token       string
	MiniAppURL  string
	Db          *gorm.DB
	BotUsername string
	mode        string
	seekerUrl   string
}

type Bot struct {
	gorm.Model
	echotron.API       `gorm:"-"`
	ChatID             int64
	Notifications      bool `gorm:"default:true"`
	StoredUsersWords   []*UsersWord
	SessionSettings    SessionOptions
	db                 *gorm.DB     `gorm:"-"`
	selfUsername       string       `gorm:"-"`
	mode               string       `gorm:"-"`
	seekerUrl          string       `gorm:"-"`
	currentWordIndex   int          `gorm:"-"`
	sessionWords       []*UsersWord `gorm:"-"`
	startSessionWords  []*UsersWord `gorm:"-"`
	sessionWordEntries []*WordEntry `gorm:"-"`
}

type SendOptions struct {
	*echotron.MessageOptions
}

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

type UsersWord struct {
	gorm.Model
	BotId               uint
	Word                string
	rememberRating      int
	sessionMistakes     uint
	usedInSession       bool
	LastSessionMistakes uint
	IsNewWord           bool `gorm:"default:true"`
	isFrontCard         bool
	reference           *UsersWord
}

type SessionOptions struct {
	BotId                     uint `gorm:"primarykey"`
	WithWordToDefinitionCards bool
	WithDefinitionToWordCards bool
	MaxDefinitionsCount       int
	MaxExamplesCount          int
}
