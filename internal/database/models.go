package database

import (
	"time"
)

type URL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OriginalURL string    `gorm:"not null" json:"originalURL"`
	ShortCode   string    `gorm:"uniqueIndex;not null" json:"shortCode"`
	UserID      *uint     `gorm:"index" json:"userId,omitempty"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Click struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URLID     uint      `gorm:"not null" json:"url_id"`
	URL       URL       `gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE" json:"url"`
	IPAddress string    `gorm:"not null" json:"ip_address"`
	Referer   *string   `json:"referer"`
	LangCode  *string   `json:"lang_code"`
	Browser   *string   `json:"browser"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
}

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	JoinedAt   time.Time `gorm:"autoCreateTime" json:"joinedAt"`
	Username   string    `gorm:"uniqueIndex;not null" json:"username"`
	Password   string    `gorm:"not null" json:"-"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	IsVerified bool      `gorm:"default:false" json:"-"`
	Urls       []URL     `gorm:"foreignKey:UserID;references:ID" json:"urls,omitempty"`
}
