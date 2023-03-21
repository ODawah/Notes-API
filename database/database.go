package database

import (
	"os"

	"github.com/Notes-App/schemas"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB.AutoMigrate(&schemas.Note{}, &schemas.User{})
	return nil
}

func CleanUp() error {
	err := os.Remove("notes.db")
	if err != nil {
		return err
	}
	return nil
}
