package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(dsn string) {
	con, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	con.AutoMigrate(&URL{}, &Click{}, &User{})
	db = con
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized")
		os.Exit(1)
	}
	return db
}
