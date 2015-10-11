package file_manager_test

import (
	"fmt"
	fm "github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	logger "github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestFileManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileManager Suite")
}

var (
	l       *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM-TEST] "})
)

var _ = BeforeSuite(func() {
	// Load test ENV variables
	godotenv.Load("../.env.test")

	fm.Logger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM] "})
})

var _ = AfterSuite(func() {
	os.RemoveAll("tmp")
})

func createTestFile(fileName string) {
	var (
		err error
		nf  *os.File
	)
	if nf, err = os.Create(fileName); err != nil {
		Fail(fmt.Sprintf("Unable to create file %s", fileName))
	}
	nf.Close()
}
