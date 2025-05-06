package database

import (
	"log"

	"gorm.io/gorm"
)

func SaveURL(db *gorm.DB, originalURL string, shortCode string) (URL, error) {
	url := URL{OriginalURL: originalURL, ShortCode: shortCode}
	if err := db.Create(&url).Error; err != nil {
		log.Println("Error when saving a url", originalURL, shortCode)
		return URL{}, err
	}
	return url, nil
}

func GetURL(db *gorm.DB, shortCode string) (string, uint, error) {
	var url URL
	if err := db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		log.Println("Error when get a url", shortCode)
		return "", 0, err
	}
	return url.OriginalURL, url.ID, nil
}

func DeleteURL(db *gorm.DB, shortCode string) error {
	var url URL
	if err := db.Where("short_code = ?", shortCode).Delete(&url).Error; err != nil {
		log.Println("Error when deleting a url", shortCode)
		return err
	}
	return nil
}

func RegisterClick(db *gorm.DB, id uint, ip string) error {
	click := Click{URLID: id, IPAddress: ip}
	if err := db.Create(&click).Error; err != nil {
		log.Println("Error when register a click", click)
		return err
	}
	return nil
}
