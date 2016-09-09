package config

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"database/sql"
)

// Return a new gorm instance with a connection to the DB as specified
// in the OS environment.
func NewDB(connection... interface{}) *gorm.DB {
	var (
		db gorm.DB
		err error
	)

	if len(connection) == 0 {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	} else {
		switch v := connection[0].(type) {
		case *sql.DB:
			db, err = gorm.Open("postgres", v)
		case *gorm.DB:
		default:
			log.Fatalln("Unable to initialize DB")
		}
	}
	// Uncomment this to enable DB logging
	db.LogMode(false)

	if err != nil {
		log.Fatalf("DB connection error: %v\n", err)
	}

	return &db
}
