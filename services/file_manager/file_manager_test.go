package file_manager_test

import (
	"fmt"
	fm "github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"time"
	"errors"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"strings"
)

var _ = Describe("FileManager", func() {

	Describe("Importing files", func() {

		Context("Having one file manager", func() {

			BeforeEach(func() {
				initFileManager()
			})

			AfterEach(func() {
				fileManager.Destroy()
				fileManager = nil
			})

			It("must prevent watching a->b and a->c simultaneously", func() {
				fileManager.Watch(watchDir1, "a->b")
				err = fileManager.Watch(watchDir1, "a->c")
				Ω(err).Should(HaveOccurred())
			})

			It("must create source and target directories if not exist", func() {
				fileManager.Watch(watchDir1, "label")

				_, err = os.Stat(watchDir1)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = os.Stat(targetDir1)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("must copy existing file from watch dir to target dir", func() {

				os.MkdirAll(watchDir1, os.ModePerm)

				createTestFile(watchFile1)

				fileManager.Watch(watchDir1, "label")

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})

			It("must copy only files, not directories", func() {
				subdir := filepath.Join(watchDir1, "subdir")
				os.MkdirAll(subdir, os.ModePerm)
				fileManager.Watch(watchDir1, targetDir1)
				time.Sleep(1 * time.Second)
				_, err = os.Stat(subdir)
				Ω(err).ShouldNot(HaveOccurred())

			})

			It("must copy new file from watch dir to target dir", func() {

				fileManager.Watch(watchDir1, targetDir1)

				createTestFile(watchFile1)

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})

			It("must copy 2 new different files to target dir", func() {

				fileManager.Watch(watchDir1, targetDir1)

				watchFile2 := filepath.Join(watchDir1, "file2.txt")
				targetFile2 := filepath.Join(targetDir1, "*/*/v01/file2.txt")

				createTestFile(watchFile1)
				createTestFile(watchFile2)

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile2)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})

		})

		Context("Having two file managers", func() {

			BeforeEach(func() {
				initFileManager()
				if fileManager2, err = fm.NewFM(targetDir2); err != nil {
					Fail(fmt.Sprintf("Unable to initialize FileManager2: %v", err))
				}

			})

			AfterEach(func() {
				fileManager.Destroy()
				fileManager2.Destroy()
				fileManager = nil
				fileManager2 = nil
			})

			It("both file managers should move files to target directories", func() {
				fileManager.Watch(watchDir1, targetDir1)
				fileManager2.Watch(watchDir2, targetDir2)
				createTestFile(watchFile1)
				createTestFile(watchFile2)

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile2v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})

			It("after destroying and recreating both file managers files should be moved to target directories", func() {
				fileManager.Destroy()
				fileManager2.Destroy()
				fileManager = nil
				fileManager2 = nil


				if fileManager, err = fm.NewFM(targetDir1, fm.WatchPair{watchDir1, "fm1"}); err != nil {
					Fail(fmt.Sprintf("Unable to initialize FileManager: %v", err))
				}
				if fileManager2, err = fm.NewFM(targetDir2, fm.WatchPair{watchDir2, "fm2"}); err != nil {
					Fail(fmt.Sprintf("Unable to initialize FileManager2: %v", err))
				}

				createTestFile(watchFile1)
				createTestFile(watchFile2)

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())

				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile2v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})
		})

		Context("copying more than once the same file to the same dir", func() {
			BeforeEach(func() {
				initFileManager()
			})

			AfterEach(func() {
				fileManager.Destroy()
				fileManager = nil
			})
			It("both files must be versioned", func() {

				fileManager.Watch(watchDir1, "label")

				createTestFile(watchFile1)


				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())

				createTestFile(watchFile1)
				Eventually(func() []string {
					files, _ := filepath.Glob(targetFile1v1)
					return files
				}, 3 * time.Second).ShouldNot(BeNil())
			})

		})
	})

	Describe("Support handlers", func() {
		BeforeEach(func() {
			initFileManager()
		})

		AfterEach(func() {
			fileManager.Destroy()
			fileManager = nil
		})

		It("registers handlers and calls them", func() {
			handlerWasCalled := false

			handler := func(file *models.File) (err error) {
				handlerWasCalled = true
				return
			}

			fileManager.Register(handler)
			fileManager.Watch(watchDir1, targetDir1)
			createTestFile(watchFile1)

			Eventually(func() bool {
				return handlerWasCalled
			}, 3 * time.Second).Should(BeTrue())
		})

		It("calls handler with proper params", func() {
			handlerWasCalled := ""

			handler := func(file *models.File) (err error) {
				handlerWasCalled = file.SourcePath + " " + file.TargetDir
				l.Println(handlerWasCalled)
				return
			}

			fileManager.Register(handler)
			fileManager.Watch(watchDir1, targetDir1)
			createTestFile(watchFile1)

			Eventually(func() bool {
				return strings.HasPrefix(handlerWasCalled, watchFile1 + " " + targetDir1)
			}, 3 * time.Second).Should(BeTrue())
		})

		It("calls more than one handler", func() {
			handlersCalled := 0

			handler1 := func(file *models.File) (err error) {
				handlersCalled ++
				return
			}

			handler2 := func(file *models.File) (err error) {
				handlersCalled ++
				return
			}

			fileManager.Register(handler1, handler2)

			fileManager.Watch(watchDir1, targetDir1)
			createTestFile(watchFile1)

			Eventually(func() bool {
				return handlersCalled == 2
			}, 3 * time.Second).Should(BeTrue())

		})
		It("calls handlers until error", func() {
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

			Consistently(func() bool {
				return handlersCalled
			}, 3 * time.Second).Should(BeFalse())

		})
	})

	Describe("Database Integrity", func() {
		BeforeEach(func() {
			initFileManager()
		})

		AfterEach(func() {
			fileManager.Destroy()
			fileManager = nil
		})

		It("must create versions for the same file and version must increment by one", func() {
			fileManager.Watch(watchDir1, "label")
			createTestFile(watchFile1)

			file := models.File{FileName: filepath.Base(watchFile1)}
			Eventually(func() error {
				return file.Load()
			}, 3 * time.Second).ShouldNot(HaveOccurred())
			Ω(file.Version).Should(Equal(int64(1)))

			createTestFile(watchFile1)
			file = models.File{FileName: filepath.Base(watchFile1), Version: 2}
			Eventually(func() error {
				return file.Load()
			}, 3 * time.Second).ShouldNot(HaveOccurred())
		})

		It("must create a file record in db", func() {
			fileManager.Watch(watchDir1, "label")

			createTestFile(watchFile1)

			//check that file is in db
			file := models.File{FileName: filepath.Base(watchFile1)}
			Eventually(func() error {
				return file.Load()
			}, 3 * time.Second).ShouldNot(HaveOccurred())


		})
	})

	Describe("Validating files", func() {
		XIt("must validate id3", func() {

		})

		XContext("When file is valid", func() {
			XIt("mark file as valid", func() {
			})
		})

		XContext("When file is invalid", func() {
			XIt("mark file as invalid", func() {
			})
			XIt("send notification to admin", func() {
			})
		})
	})
})
