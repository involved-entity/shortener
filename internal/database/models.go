package database

import (
	"time"
)

type URL struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	ShortCode   string `gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Click struct {
	ID        uint      `gorm:"primaryKey"`
	URLID     uint      `gorm:"not null"`
	URL       URL       `gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE"`
	IPAddress string    `gorm:"not null"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
}
