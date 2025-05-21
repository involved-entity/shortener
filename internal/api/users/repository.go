package users

import (
	"log"
	"shortener/internal/database"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type UserInfo struct {
	Username string
	ID       int
}

func (r Repository) SaveUser(username string, email string, password string) (database.User, error) {
	user := database.User{Username: username, Email: email, Password: password}
	if err := r.db.Create(&user).Error; err != nil {
		log.Println("Error when saving user", user)
		return database.User{}, err
	}
	return user, nil
}

func (r Repository) GetUser(userInfo UserInfo) (database.User, error) {
	var user database.User
	var err error
	if userInfo.ID != 0 {
		err = r.db.Where("id = ? AND is_verified = true", userInfo.ID).First(&user).Error
	} else {
		err = r.db.Where("username = ? AND is_verified = true", userInfo.Username).First(&user).Error
	}
	if err != nil {
		log.Println("Error when get a user", userInfo, err)
		return database.User{}, err
	}
	return user, nil
}

func (r Repository) VerificateUser(id int) error {
	var user database.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println("Error when get a user", id, err)
		return err
	}
	user.IsVerified = true
	if err := r.db.Save(&user).Error; err != nil {
		log.Println("Error when save user verified status", err)
		return err
	}
	return nil
}

func (r Repository) ChangeUserPassword(id int, hashedPassword string) error {
	if err := r.db.Model(&database.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		log.Println("Error when set new password for user", id, err)
		return err
	}
	return nil
}
