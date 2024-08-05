package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbHost string

func ConnectDatabase() (*gorm.DB, error) {
	dbHost = "postgres"
	log.Println("Starting db postgres!")

	dsn := fmt.Sprintf("host=%s user=postgres password=changeme dbname=gorm port=5432 sslmode=disable TimeZone=Europe/Minsk", dbHost)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Bot{}, &WordEntry{}, &SpeechPartEntry{}, &WordDefinition{}, &WordUsageExample{}, &SentencePice{})

	//db.Create(&Product{Code: "D42", Price: 100})
	return db, nil
}

// TODO change to db.func
func saveBotToDbIfNotExists(db *gorm.DB, b *Bot) {
	var tryBot Bot
	db.FirstOrCreate(
		&tryBot, Bot{ChatID: b.ChatID},
	)
}
