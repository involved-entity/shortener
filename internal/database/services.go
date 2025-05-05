package database

import (
	"log"

	"gorm.io/gorm"
)

func SaveURL(db *gorm.DB, originalURL string, shortCode string) {
	url := URL{OriginalURL: originalURL, ShortCode: shortCode}
	if err := db.Create(&url); err != nil {
		log.Println("Error when saving a url", originalURL, shortCode)
	}
	log.Println(url.ID)
}
