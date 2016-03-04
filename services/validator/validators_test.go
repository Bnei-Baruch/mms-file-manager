package validator_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"database/sql"
	"github.com/Bnei-Baruch/mms-file-manager/utils"
	v "github.com/Bnei-Baruch/mms-file-manager/services/validator"
	"time"
	"fmt"
)

func TestValidators(t *testing.T) {
	db := utils.SetupSpec()
	fmt.Println("TestValidators start time", time.Now())
	Convey("Workflow with list of validations", t, func() {
		db.Exec("DELETE FROM patterns; DELETE FROM files; DELETE FROM workflows;")
		pattern := &models.Pattern{
			Name: "lang_arutz_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arutz"},
				{Key: "date", },
				{Key: "content_type", },
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		}
		pattern.Save()

		workflow := &models.Workflow{
			PatternId: sql.NullInt64{Int64: pattern.ID, Valid: true},
			EntryPoint: "SiumAvoda",
			ContentType: sql.NullString{String:"*", Valid: true},
			Line: sql.NullString{String:"*", Valid: true},
		}
		workflow.Save()

		file := &models.File{
			FileName: "heb_arutz_2012-12-16_film_crossroads.mpg",
			TargetDir: "targetDir",
			EntryPoint: "SiumAvoda",
			SourcePath: "path",
			FullPath: "../../test_files/heb_o_rav_achana_2015-10-13_lesson.mp4",
			PatternId: sql.NullInt64{Int64: pattern.ID, Valid: true},
			WorkflowId: sql.NullInt64{Int64: workflow.ID, Valid: true},
			Status: models.HAS_WORKFLOW,
			Attributes: models.JSONB{"date": "2012-12-16", "lang": "heb", "line": "crossroads", "archive_type": "arutz", "content_type": "film"},

		}

		file.CreateVersion()

		Convey("When there is an existing validation", func() {
			workflow.Validations = models.StringSlice{"passedValidation", "failedValidation", "checkFrameRate"}
			workflow.Save()
			Convey("It should run all workflow's validations", func() {
				err := v.RunValidations(file)
				So(err, ShouldBeNil)
				So(file.ValidationResult, ShouldNotBeNil)

				/*
				validationResult := models.JSONB{
					"passedValidation": models.ValidationResult{Passed: true, },
					"failedValidation": models.ValidationResult{Passed: false, ErrorMessage: "error"},
				}
				*/

				//TODO: So(validationCount(file.ValidationResult), ShouldEqual, len(workflow.GetValidations()))
				//TODO: Test that validation names are equal to...
			})
		})
		/*
		Convey("When there is a non existing validation", func() {
			Convey("It should be removed from workflow", func() {
			})
		})
		Convey("when running kuku validation", func() {
			Convey("It should return expected result", func() {
			})
		})
		Convey("When all validations pass", func() {

			Convey("File should have status VALIDATION_PASSED and list of validation results", func() {

			})

		})
		Convey("When there are no validations", func() {

			Convey("File should have status VALIDATION_PASSED and list of validation results should be empty", func() {

			})

		})
		Convey("When not all validations pass", func() {

			Convey("File should have status VALIDATION_FAILED and list of validation results", func() {

			})

		})
		*/
	})
	fmt.Println("TestValidators end time", time.Now())
}

func validationCount() {}