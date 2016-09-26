package file_manager_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	fm "github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/Bnei-Baruch/mms-file-manager/test_helpers"
	"github.com/Bnei-Baruch/mms-file-manager/utils"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	l                     *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[FM-TEST] "})
	db                    *gorm.DB
	watchDir1, targetDir1 = "tmp/source1", "tmp/target1"
	watchFile1            = filepath.Join(watchDir1, "file1.txt")
	targetFile1v1         = filepath.Join(targetDir1, "*/*/v01/file1.txt")
	targetFile1v2         = filepath.Join(targetDir1, "*/*/v02/file1.txt")

	watchDir2, targetDir2 = "tmp/source2", "tmp/target2"
	watchFile2            = filepath.Join(watchDir2, "file2.txt")
	targetFile2v1         = filepath.Join(targetDir2, "*/*/v01/file2.txt")

	fileManager  *fm.FileManager = nil
	fileManager2 *fm.FileManager = nil
	err          error
)

func TestMain(m *testing.M) {

	db = test_helpers.SetupSpec()
	/*
		flag.Parse()
		os.RemoveAll("tmp")
		// Load test ENV variables
		godotenv.Load("../../.env.test")
		db = config.NewDB()
		models.New(db)
		if err := gofer.LoadAndPerform("db:empty", "--env=../../.env.test"); err != nil {
			panic(fmt.Sprintf("Unable to empty database %v", err))
		}

		if err := gofer.LoadAndPerform("db:migrate", "--env=../../.env.test"); err != nil {
			panic(fmt.Sprintf("Unable to migrate database %v", err))
		}
	*/
	os.Exit(m.Run())
}

func TestFileManagerSpec(t *testing.T) {
	Convey("FileManager", t, func() {

		Convey("Importing files", func() {

			Convey("Having one file manager", func() {

				initFileManager()

				Reset(func() {
					fileManager.Destroy()
					fileManager = nil
				})

				Convey("It should prevent watching a->b and a->c simultaneously", func() {
					fileManager.Watch(watchDir1, "a->b")
					err = fileManager.Watch(watchDir1, "a->c")
					So(err, ShouldNotBeNil)
				})

				Convey("It should create source and target directories if not exist", func() {
					fileManager.Watch(watchDir1, "label")

					_, err = os.Stat(watchDir1)
					So(err, ShouldBeNil)

					_, err = os.Stat(targetDir1)
					So(err, ShouldBeNil)
				})

				Convey("It should copy existing file from watch dir to target dir", func() {

					os.MkdirAll(watchDir1, os.ModePerm)

					createTestFile(watchFile1)

					fileManager.Watch(watchDir1, "label")
					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})

				Convey("It should copy only files, not directories", func() {
					subdir := filepath.Join(watchDir1, "subdir")
					os.MkdirAll(subdir, os.ModePerm)
					fileManager.Watch(watchDir1, targetDir1)
					time.Sleep(1 * time.Second)
					_, err = os.Stat(subdir)
					So(err, ShouldBeNil)

				})

				Convey("It should copy new file from watch dir to target dir", func() {

					fileManager.Watch(watchDir1, targetDir1)

					createTestFile(watchFile1)
					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})

				Convey("It should copy 2 new different files to target dir", func() {

					fileManager.Watch(watchDir1, targetDir1)

					watchFile2 := filepath.Join(watchDir1, "file2.txt")
					targetFile2 := filepath.Join(targetDir1, "*/*/v01/file2.txt")

					createTestFile(watchFile1)
					createTestFile(watchFile2)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile2)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})

			})

			Convey("Having two file managers", func() {

				initFileManager()
				if fileManager2, err = fm.NewFM(targetDir2); err != nil {
					panic(fmt.Sprintf("Unable to initialize FileManager2: %v", err))
				}

				Reset(func() {
					fileManager.Destroy()
					fileManager2.Destroy()
					fileManager = nil
					fileManager2 = nil
				})

				Convey("both file managers should move files to target directories", func() {
					fileManager.Watch(watchDir1, targetDir1)
					fileManager2.Watch(watchDir2, targetDir2)
					createTestFile(watchFile1)
					createTestFile(watchFile2)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile2v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})

				Convey("after destroying and recreating both file managers files should be moved to target directories", func() {
					fileManager.Destroy()
					fileManager2.Destroy()
					fileManager = nil
					fileManager2 = nil

					if fileManager, err = fm.NewFM(targetDir1, fm.WatchPair{watchDir1, "fm1"}); err != nil {
						panic(fmt.Sprintf("Unable to initialize FileManager: %v", err))
					}
					if fileManager2, err = fm.NewFM(targetDir2, fm.WatchPair{watchDir2, "fm2"}); err != nil {
						panic(fmt.Sprintf("Unable to initialize FileManager2: %v", err))
					}

					createTestFile(watchFile1)
					createTestFile(watchFile2)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile2v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})
			})

			Convey("copying more than once the same file to the same dir", func() {
				initFileManager()

				Reset(func() {
					fileManager.Destroy()
					fileManager = nil
				})
				Convey("both files must be versioned", func() {

					fileManager.Watch(watchDir1, "label")

					createTestFile(watchFile1)

					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)

					createTestFile(watchFile1)
					So(func() interface{} {
						files, _ := filepath.Glob(targetFile1v1)
						return files
					}, utils.Eventually, 3*time.Second, assertions.ShouldNotBeNil)
				})

			})
		})

		Convey("Support handlers", func() {
			initFileManager()

			Reset(func() {
				fileManager.Destroy()
				fileManager = nil
			})

			Convey("It registers handlers and calls them", func() {
				handlerWasCalled := false

				handler := func(file *models.File) (err error) {
					handlerWasCalled = true
					return
				}

				fileManager.Register(handler)
				fileManager.Watch(watchDir1, targetDir1)
				createTestFile(watchFile1)

				So(func() interface{} {
					return handlerWasCalled
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeTrue)
			})

			Convey("It calls handler with proper params", func() {
				handlerWasCalled := ""

				handler := func(file *models.File) (err error) {
					handlerWasCalled = file.SourcePath + " " + file.TargetDir
					l.Println(handlerWasCalled)
					return
				}

				fileManager.Register(handler)
				fileManager.Watch(watchDir1, targetDir1)
				createTestFile(watchFile1)

				So(func() interface{} {
					return strings.HasPrefix(handlerWasCalled, watchFile1+" "+targetDir1)
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeTrue)
			})

			Convey("calls more than one handler", func() {
				handlersCalled := 0

				handler1 := func(file *models.File) (err error) {
					handlersCalled++
					return
				}

				handler2 := func(file *models.File) (err error) {
					handlersCalled++
					return
				}

				fileManager.Register(handler1, handler2)

				fileManager.Watch(watchDir1, targetDir1)
				createTestFile(watchFile1)

				So(func() interface{} {
					return handlersCalled == 2
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeTrue)
			})
			Convey("calls handlers until error", func() {
				handlersCalled := false

				handler1 := func(file *models.File) (err error) {
					return errors.New("new error")
				}

				handler2 := func(file *models.File) (err error) {
					handlersCalled = true
					return
				}

				fileManager.Register(handler1, handler2)
				fileManager.Watch(watchDir1, targetDir1)
				createTestFile(watchFile1)

				So(func() interface{} {
					return handlersCalled
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeFalse)
			})
		})

		Convey("Database Integrity", func() {
			initFileManager()

			Reset(func() {
				fileManager.Destroy()
				fileManager = nil
			})

			Convey("It should create versions for the same file and version must increment by one", func() {
				fileManager.Watch(watchDir1, "label")
				createTestFile(watchFile1)

				file := models.File{FileName: filepath.Base(watchFile1)}

				So(func() interface{} {
					return file.Load()
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeNil)

				So(file.Version, ShouldEqual, int64(1))

				files, _ := filepath.Glob(targetFile1v1)
				So(files, ShouldNotBeNil)

				createTestFile(watchFile1)
				file = models.File{FileName: filepath.Base(watchFile1), Version: 2}

				So(func() interface{} {
					return file.Load()
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeNil)

				files, _ = filepath.Glob(targetFile1v2)
				So(files, ShouldNotBeNil)
			})

			Convey("It should create a file record in db", func() {
				fileManager.Watch(watchDir1, "label")

				createTestFile(watchFile1)

				//check that file is in db
				file := models.File{FileName: filepath.Base(watchFile1)}

				So(func() interface{} {
					return file.Load()
				}, utils.Eventually, 3*time.Second, assertions.ShouldBeNil)
			})
		})

		SkipConvey("Validating files", func() {
			Convey("It should validate id3", func() {

			})

			Convey("When file is valid", func() {
				Convey("mark file as valid", func() {
				})
			})

			Convey("When file is invalid", func() {
				Convey("mark file as invalid", func() {
				})
				Convey("send notification to admin", func() {
				})
			})
		})
	})
}

func createTestFile(fileName string) {
	var (
		err error
		nf  *os.File
	)
	if nf, err = os.Create(fileName); err != nil {
		panic(fmt.Sprintf("Unable to create file %s", fileName))
	}
	nf.Close()
}

func initFileManager() {
	db.Exec("DELETE FROM files;")

	if err = os.RemoveAll(watchDir1); err != nil {
		panic("Unable to remove watch dir")
	}

	if err = os.RemoveAll(targetDir1); err != nil {
		panic("Unable to remove target dir")
	}
	if err = os.RemoveAll(watchDir2); err != nil {
		panic("Unable to remove watch dir")
	}

	if err = os.RemoveAll(targetDir2); err != nil {
		panic("Unable to remove target dir")
	}

	if fileManager, err = fm.NewFM(targetDir1); err != nil {
		panic(fmt.Sprintf("Unable to initialize FileManager: %v", err))
	}
}
