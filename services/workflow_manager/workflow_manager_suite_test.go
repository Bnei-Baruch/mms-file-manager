package workflow_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/jinzhu/gorm"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/joho/godotenv"
	"github.com/chuckpreslar/gofer"
	"fmt"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"log"
)

var (
	l      *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[WM-TEST] "})
	db    *gorm.DB
	err error
)

func TestWorkflowManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkflowManager Suite")
}

var _ = BeforeSuite(func() {
	// Load test ENV variables
	godotenv.Load("../../.env.test")
	db = config.NewDB()
	models.New(db)

	if err := gofer.LoadAndPerform("db:empty", "--env=../../.env.test"); err != nil {
		Fail(fmt.Sprintf("Unable to empty database %v", err))
	}

	if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
		Fail(fmt.Sprintf("Unable to migrate database %v", err))
	}
})

var _ = AfterSuite(func() {
})

