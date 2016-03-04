package validator

import (
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"log"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"errors"
	"github.com/fzakaria/goav/avformat"
	"fmt"
)

var l *log.Logger = nil

func init() {
	l = logger.InitLogger(&logger.LogParams{LogPrefix: "[VAL] "})
}

type validationFunc func(f *models.File) (passed bool, err error)

var validations = map[string]validationFunc{
	"passedValidation": passedValidation,
	"failedValidation": failedValidation,
	"checkFrameRate": checkFrameRate,
}

func passedValidation(f *models.File) (passed bool, err error) {
	passed = true
	return
}

func failedValidation(f *models.File) (passed bool, err error) {
	passed = false
	err = errors.New("error")
	return
}

func checkFrameRate(f *models.File) (paased bool, err error) {
	var (
		ctxtFormat        *avformat.Context
		url string
	)

	// Register all formats and codecs
	avformat.AvRegisterAll()

	// Open video file
	if avformat.AvformatOpenInput(&ctxtFormat, f.SourcePath, nil, nil) != 0 {
		log.Println("Error: Couldn't open file.")
		return
	}

	// Retrieve stream information
	if ctxtFormat.AvformatFindStreamInfo(nil) < 0 {
		log.Println("Error: Couldn't find stream information.")
		return
	}

	// Dump information about file onto standard error
	ctxtFormat.AvDumpFormat(0, url, 0)

	return
}

var RunValidations = file_manager.HandlerFunc(func(file *models.File) (err error) {

	defer func() {
		if err != nil {
			return
		}
		if err = file.Save(); err != nil {
			l.Printf("Problem updating file -  %s : %s", file.FileName, err)
		}
	}()
	//TODO: workflow is not loaded!!!
	file.Load()
	workflow := file.Workflow
	fmt.Println("jyjy", workflow.Validations, file.Workflow.Validations, file.WorkflowId)
	for _, validation := range workflow.Validations {
		fmt.Println("kuku validations", validation)
		fn := validations[validation]
		passed, err := fn(file)
		file.ValidationResult[validation] = models.ValidationResult{Passed: passed, ErrorMessage: err}
	}

	return

})
