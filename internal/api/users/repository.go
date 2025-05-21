package users

import (
	"log"
	"shortener/internal/database"

	"gorm.io/gorm"
)

type UsersRepository interface {
	SaveUser(username string, email string, password string) (database.User, error)
}

type Repository struct {
	db *gorm.DB
}

func (r Repository) SaveUser(username string, email string, password string) (database.User, error) {
	user := database.User{Username: username, Email: email, Password: password}
	if err := r.db.Create(&user).Error; err != nil {
		log.Println("Error when saving user", user)
		return database.User{}, err
	}
	return user, nil
}

func (r Repository) GetUser(username string) (database.User, error) {
	var user database.User
	if err := r.db.Where("username = ? AND is_verified = true", username).First(&user).Error; err != nil {
		log.Println("Error when get a user", username, err)
		return database.User{}, err
	}
	return user, nil
}
