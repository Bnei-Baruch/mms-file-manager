package models
import (
	"time"
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//	DeletedAt *time.Time `sql:"index"`
}

var db *gorm.DB

func New(dbConn *gorm.DB) {
	db = dbConn
}