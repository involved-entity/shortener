package database

import (
	"time"
)

type URL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OriginalURL string    `gorm:"not null" json:"originalURL"`
	ShortCode   string    `gorm:"uniqueIndex;not null" json:"shortCode"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Click struct {
	ID        uint      `gorm:"primaryKey"`
	URLID     uint      `gorm:"not null"`
	URL       URL       `gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE"`
	IPAddress string    `gorm:"not null"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
}

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	JoinedAt   time.Time `gorm:"autoCreateTime" json:"joinedAt"`
	Username   string    `gorm:"uniqueIndex;not null" json:"username"`
	Password   string    `gorm:"not null" json:"-"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	IsVerified bool      `gorm:"default:false" json:"-"`
}
