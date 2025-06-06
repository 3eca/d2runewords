package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// connects to the database
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error while connect to sqlite:", err)
		return nil, err
	}
	return db, err
}