package utils

import (
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/chuckpreslar/gofer"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/joho/godotenv"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func SetupSpec() (db *gorm.DB) {
	fmt.Println("start time", time.Now())

	// Load test ENV variables
	godotenv.Load("../../.env.test")
	fmt.Println("1 time", time.Now())
	if err := gofer.LoadAndPerform("db:empty", "--env=../../.env.test"); err != nil {
		panic(fmt.Sprintf("Unable to empty database %v", err))
	}
	/*
	fmt.Println("2 time", time.Now())

	if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
		panic(fmt.Sprintf("Unable to migrate database %v", err))
	}
	fmt.Println("3 time", time.Now())
	*/
	db = config.NewDB()
	fmt.Println("4 time", time.Now())
	models.New(db)
	fmt.Println("end time", time.Now())
	return
}
