package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"log"
)

// Return a new gorm instance with a connection to the DB as specified
// in the OS environment.
func NewDB() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	// Uncomment this to enable DB loggin
	db.LogMode(true)

	if err != nil {
		log.Fatalf("DB connection error: %v\n", err)
	}

	return &db
}