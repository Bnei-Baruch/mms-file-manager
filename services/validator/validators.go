package validator

import (
	"github.com/go-errors/errors"

	"os/exec"

	"encoding/json"
	"fmt"

	"reflect"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/sthorne/reflections"
)

var l = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[VAL] "})

type (
	validationFunc func(f *models.File) (passed bool, err error)
)

var validations = map[string]validationFunc{
	"passedValidation": passedValidation,
	"failedValidation": failedValidation,
	"checkExif":        checkExif,
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

func checkExif(f *models.File) (passed bool, err error) {

	var cmdOut []byte
	args := []string{"-j", f.FullPath}
	if cmdOut, err = exec.Command("exiftool", args...).Output(); err != nil {
		return false, errors.Errorf("There was an error retrieving file %s metadata: %s", f.FullPath, err)
	}

	var exifs []models.Exif
	if err = json.Unmarshal(cmdOut, &exifs); err != nil {
		return false, errors.Errorf("There was an error parsing file %s metadata: %s", f.FullPath, err)
	}
	f.Exif = exifs[0]
	f.Save()
	workflowExif, _ := reflections.Items(f.Workflow.Exif)

	var checkError string
	for fieldName, fieldValue := range workflowExif {
		if isZeroOfUnderlyingType(fieldValue) {
			continue
		}
		var value interface{}
		if value, err = reflections.GetField(f.Exif, fieldName); err != nil {
			return false, errors.Errorf("There was an error validating Exif for file %s. Field '%s' doesn't exist:\n %s", f.FullPath, fieldName, err)
		}
		if value != fieldValue {
			checkError += fmt.Sprintf("Field '%s' has value '%v', but '%v' was expected\n", fieldName, value, fieldValue)
		}

	}

	if checkError != "" {
		return false, errors.Errorf("There was an error validating Exif for file %s.\n%s", f.FullPath, checkError)
	}
	return true, nil

}

// RunValidations will run all registeded validations
var RunValidations = file_manager.HandlerFunc(func(file *models.File) (err error) {

	defer func() {
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

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
