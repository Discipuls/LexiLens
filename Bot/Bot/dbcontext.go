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
	err = db.AutoMigrate(&Bot{}, &WordEntry{}, &SpeechPartEntry{}, &WordDefinition{}, &DefinitionPiece{}, &WordUsageExample{}, &SentencePice{}, &UsersWord{},
		&SessionOptions{})
	if err != nil {
		log.Println("Error: db.Automigrate: ", err.Error())
	}

	return db, nil
}

func insertBotToDbIfNotExists(db *gorm.DB, chatID int64) *Bot {
	var tryBot Bot
	tx := db.Preload("StoredUsersWords").Preload("SessionSettings").FirstOrCreate(
		&tryBot, Bot{ChatID: chatID},
	)
	if tx.Error != nil {
		log.Println("insertBotToDbIfNotExists error: ", tx.Error.Error())
	}
	return &tryBot
}

func insertSessionItemToDbIfNotExists(db *gorm.DB, item UsersWord) *UsersWord {
	var dbItem UsersWord
	db.FirstOrCreate(
		&dbItem, item,
	)
	return &dbItem
}

func SaveBotWordEntryInDb(db *gorm.DB, item *UsersWord) {
	tx := db.Save(*item)
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

func SaveSessionOptionsToDb(db *gorm.DB, opts *SessionOptions) {
	Wg.Add(1)
	go SaveSessionOptionsToDbWg(db, opts)
}

func SaveSessionOptionsToDbWg(db *gorm.DB, opts *SessionOptions) {
	defer Wg.Done()
	tx := db.Save(&opts)
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
