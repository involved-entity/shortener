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

func GetURL(db *gorm.DB, shortCode string) (string, error) {
	var url URL
	if err := db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		log.Println("Error when get a url", shortCode)
		return "", err
	}
	return url.OriginalURL, nil
}
