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
	"github.com/chuckpreslar/gofer"
	"github.com/Bnei-Baruch/mms-file-manager/config"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/jinzhu/gorm"
	"path/filepath"
)

func TestFileManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileManager Suite")
}

var (
	l      *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM-TEST] "})
	db    *gorm.DB
	watchDir1, targetDir1 = "tmp/source1", "tmp/target1"
	watchFile1 = filepath.Join(watchDir1, "file1.txt")
	targetFile1v1 = filepath.Join(targetDir1, "*/*/v01/file1.txt")
	targetFile1v2 = filepath.Join(targetDir1, "*/*/v02/file1.txt")

	watchDir2, targetDir2 = "tmp/source2", "tmp/target2"
	watchFile2 = filepath.Join(watchDir2, "file2.txt")
	targetFile2v1 = filepath.Join(targetDir2, "*/*/v01/file2.txt")

	fileManager  *fm.FileManager = nil
	fileManager2 *fm.FileManager = nil
	err error

)

var _ = BeforeSuite(func() {
	os.RemoveAll("tmp")
	// Load test ENV variables
	godotenv.Load("../../.env.test")
	if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
		Fail(fmt.Sprintf("Unable to migrate database %v", err))
	}
	db = config.NewDB()
	models.New(db)
	fm.Logger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM] "})
})

var _ = AfterSuite(func() {
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

func initFileManager() {
	db.Exec("DELETE FROM files;")

	if err = os.RemoveAll(watchDir1); err != nil {
		Fail("Unable to remove watch dir")
	}

	if err = os.RemoveAll(targetDir1); err != nil {
		Fail("Unable to remove target dir")
	}
	if err = os.RemoveAll(watchDir2); err != nil {
		Fail("Unable to remove watch dir")
	}

	if err = os.RemoveAll(targetDir2); err != nil {
		Fail("Unable to remove target dir")
	}

	if fileManager, err = fm.NewFM(targetDir1); err != nil {
		Fail(fmt.Sprintf("Unable to initialize FileManager: %v", err))
	}
}
