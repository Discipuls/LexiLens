package bot

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(conf *config) (*gorm.DB, error) {
	log.Println("Connection to postgres database...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.DatabaseName,
		conf.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error opening database: ", err.Error())
	}
	err = db.AutoMigrate(&Bot{}, &WordEntry{}, &SpeechPartEntry{}, &WordDefinition{}, &DefinitionPiece{}, &WordUsageExample{}, &SentencePice{}, &BotWordEntry{})
	if err != nil {
		log.Println("Error: db.Automigrate: ", err.Error())
	}

	return db, nil
}

func insertBotToDbIfNotExists(db *gorm.DB, chatID int64) *Bot {
	var tryBot Bot
	db.Preload("WordEntries").FirstOrCreate(
		&tryBot, Bot{ChatID: chatID},
	)
	return &tryBot
}

func insertBotWordEntryToDbIfNotExists(db *gorm.DB, entry BotWordEntry) *BotWordEntry {
	var tryEntry BotWordEntry
	db.FirstOrCreate(
		&tryEntry, entry,
	)
	return &tryEntry
}

func SaveBotWordEntryInDb(db *gorm.DB, entry *BotWordEntry) {
	tx := db.Save(*entry)
	if tx.Error != nil {
		log.Println("SaveBotWordEntryInDb error: ", tx.Error.Error())
	}
}

func FindWordEntry(db *gorm.DB, word string) (*WordEntry, error) {
	var tryEntry WordEntry
	tx := db.Preload("SpeechParts").
		Preload("SpeechParts.Definitions").
		Preload("SpeechParts.Definitions.Definition").
		Preload("SpeechParts.Definitions.WordUsageExamples").
		Preload("SpeechParts.Definitions.WordUsageExamples.Pieces").First(
		&tryEntry, WordEntry{Word: word},
	)
	return &tryEntry, tx.Error

}

func InsertWordEntryToDb(db *gorm.DB, wordEntry *WordEntry) error {
	tx := db.Save(wordEntry)
	return tx.Error
}

func SaveBotToDb(db *gorm.DB, b *Bot) {
	Wg.Add(1)
	go SaveBotToDbWg(db, b)
}

func SaveBotToDbWg(db *gorm.DB, b *Bot) {
	defer Wg.Done()
	tx := db.Save(&b)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
}

func SaveUsersWordToDb(db *gorm.DB, wordEntry *WordEntry) {
	Wg.Add(1)
	saveUsersWordToDbWg(db, wordEntry)
}

func saveUsersWordToDbWg(db *gorm.DB, wordEntry *WordEntry) {
	defer Wg.Done()
	tx := db.Save(&wordEntry)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}

}
