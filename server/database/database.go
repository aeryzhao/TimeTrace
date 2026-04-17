package database

import (
	"log"
	"os"
	"timetrace/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dbPath := os.Getenv("TIMETRACE_DB_PATH")
	if dbPath == "" {
		dbPath = "timetrace.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	err = DB.AutoMigrate(&models.Category{}, &models.Activity{}, &models.TimeEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
