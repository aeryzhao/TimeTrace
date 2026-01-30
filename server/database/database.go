package database

import (
	"log"
	"timetrace/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("timetrace.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	err = DB.AutoMigrate(&models.Category{}, &models.Activity{}, &models.TimeEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
