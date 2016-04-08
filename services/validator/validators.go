package validator

import (
	"errors"
	"log"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/fzakaria/goav/avformat"
)

var l *log.Logger = nil

func init() {
	l = logger.InitLogger(&logger.LogParams{LogPrefix: "[VAL] "})
}

type validationFunc func(f *models.File) (passed bool, err error)

var validations = map[string]validationFunc{
	"passedValidation": passedValidation,
	"failedValidation": failedValidation,
	"checkFrameRate":   checkFrameRate,
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

func checkFrameRate(f *models.File) (passed bool, err error) {
	var (
		ctxtFormat *avformat.Context
		url        string
	)

	// Register all formats and codecs
	avformat.AvRegisterAll()

	// Open video file
	if avformat.AvformatOpenInput(&ctxtFormat, f.SourcePath, nil, nil) != 0 {
		return false, errors.New("checkFrameRate Error: Couldn't open file.")
	}

	// Retrieve stream information
	if ctxtFormat.AvformatFindStreamInfo(nil) < 0 {
		return false, errors.New("checkFrameRate Error: Couldn't find stream information.")
	}

	// Dump information about file onto standard error
	ctxtFormat.AvDumpFormat(0, url, 0)

	return true, nil
}

var RunValidations = file_manager.HandlerFunc(func(file *models.File) (err error) {

	defer func() {
		if err != nil {
			return
		}
		if err = file.Save(); err != nil {
			l.Printf("RunValidations Error: Problem updating file -  %s : %s", file.FileName, err)
		}
	}()
	file.Load()
	workflow := file.Workflow
	for _, validation := range workflow.Validations {
		fn := validations[validation]
		passed, err := fn(file)
		var msg string
		if err != nil {
			msg = err.Error()
		}
		file.ValidationResult[validation] = models.ValidationResult{Passed: passed, ErrorMessage: msg}
	}

	return

})
