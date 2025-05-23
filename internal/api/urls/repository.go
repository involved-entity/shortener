package urls

import (
	"log"

	conf "shortener/internal/config"
	"shortener/internal/database"

	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	UserId int
	Page   int
}

func (r Repository) SaveURL(originalURL string, shortCode string) (database.URL, error) {
	var userID *uint
	if r.UserId != 0 {
		id := uint(r.UserId)
		userID = &id
	}
	url := database.URL{OriginalURL: originalURL, ShortCode: shortCode, UserID: userID}
	if err := r.db.Create(&url).Error; err != nil {
		log.Println("Error when saving a url", originalURL, shortCode, userID)
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
		log.Println("Error when deleting a url", shortCode, r.UserId)
		return err
	}
	return nil
}

func (r Repository) RegisterClick(id uint, ip string, referer string, langCode string, browser string) error {
	click := database.Click{URLID: id, IPAddress: ip, Referer: &referer, LangCode: &langCode, Browser: &browser}
	if err := r.db.Create(&click).Error; err != nil {
		log.Println("Error when register a click", click)
		return err
	}
	return nil
}

func (r Repository) GetUserURLs() ([]database.URL, error) {
	var urls []database.URL
	limit := conf.GetConfig().PageSize
	if err := r.db.
		Model(&database.URL{}).
		Joins("User").
		Offset((r.Page-1)*limit).
		Limit(limit).
		Where("user_id = ?", r.UserId).
		Find(&urls).
		Error; err != nil {
		log.Println("Error when get user urls", err)
		return urls, err
	}
	return urls, nil
}

func (r Repository) GetURLClicks(shortCode string) ([]database.Click, error) {
	var clicks []database.Click
	limit := conf.GetConfig().PageSize
	if err := r.db.
		Model(&database.Click{}).
		Joins("URL").
		Joins("URL.User").
		Offset((r.Page-1)*limit).
		Limit(limit).
		Where(`"URL"."short_code" = ?`, shortCode).
		Where(`"URL"."user_id" = ?`, r.UserId).
		Find(&clicks).
		Error; err != nil {
		log.Println("Error when get user urls", err)
		return clicks, err
	}
	return clicks, nil
}
