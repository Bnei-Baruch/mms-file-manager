package event_manager_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/joho/godotenv"
	"github.com/chuckpreslar/gofer"
	"fmt"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	l      *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[EM-TEST] "})
	db    *gorm.DB
	err error
)

func TestEventManagerSpec(t *testing.T) {
	// Load test ENV variables
	godotenv.Load("../../.env.test")
	if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
		panic(fmt.Sprintf("Unable to migrate database %v", err))
	}
	db = config.NewDB()
	models.New(db)

SkipConvey("EventManager", t, func() {
	Convey("Creating new event", func() {
		Convey("Checks if file name is valid", func() {  })
		Convey("", func() {  })
		Convey("When file name is valid", func() {
			Convey("It must create new event with status NEW", func() {  })
			Convey("When event with this file already exists", func() {
				Convey("It must add/replace file in event", func() {  })
			})
		})

		Convey("When file name is invalid", func() {

		})
	})
})
}
