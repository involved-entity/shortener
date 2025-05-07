package urls

import (
	"log"

	"shortener/internal/database"

	"gorm.io/gorm"
)

type URLRepository interface {
	SaveURL(originalURL string, shortCode string) (database.URL, error)
	GetURL(shortCode string) (string, uint, error)
	DeleteURL(shortCode string) error
	RegisterClick(id uint, ip int) error
}

type Repository struct {
	db     *gorm.DB
	UserId int
}

func (r Repository) SaveURL(originalURL string, shortCode string) (database.URL, error) {
	var userID *uint
	if r.UserId != 0 {
		id := uint(r.UserId)
		userID = &id
	}
	url := database.URL{OriginalURL: originalURL, ShortCode: shortCode, UserID: userID}
	if err := r.db.Create(&url).Error; err != nil {
		log.Println("Error when saving a url", originalURL, shortCode)
		return database.URL{}, err
	}
	return url, nil
}

func (r Repository) GetURL(shortCode string) (string, uint, error) {
	var url database.URL
	if err := r.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		log.Println("Error when get a url", shortCode)
		return "", 0, err
	}
	return url.OriginalURL, url.ID, nil
}

func (r Repository) DeleteURL(shortCode string) error {
	var url database.URL
	if err := r.db.Where("short_code = ?", shortCode).Where("user_id = ?", r.UserId).Delete(&url).Error; err != nil {
		log.Println("Error when deleting a url", shortCode)
		return err
	}
	return nil
}

func (r Repository) RegisterClick(id uint, ip string) error {
	click := database.Click{URLID: id, IPAddress: ip}
	if err := r.db.Create(&click).Error; err != nil {
		log.Println("Error when register a click", click)
		return err
	}
	return nil
}
