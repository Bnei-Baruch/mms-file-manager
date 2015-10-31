package event_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"log"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/joho/godotenv"
	"github.com/chuckpreslar/gofer"
	"fmt"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/jinzhu/gorm"
	"github.com/Bnei-Baruch/mms-file-manager/models"
)

var (
	l      *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[EM-TEST] "})
	db    *gorm.DB
	err error
)

func TestEventManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EventManager Suite")
}

var _ = BeforeSuite(func() {
	// Load test ENV variables
	godotenv.Load("../../.env.test")
	if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
		Fail(fmt.Sprintf("Unable to migrate database %v", err))
	}
	db = config.NewDB()
	models.New(db)
//	fm.Logger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM] "})
})

var _ = AfterSuite(func() {
})
