package validator_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	v "github.com/Bnei-Baruch/mms-file-manager/services/validator"
	"github.com/Bnei-Baruch/mms-file-manager/test_helpers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidators(t *testing.T) {
	db := test_helpers.SetupSpec()
	fmt.Println("TestValidators start time", time.Now())
	Convey("Workflow with list of validations", t, func() {
		db.Exec("DELETE FROM patterns; DELETE FROM files; DELETE FROM workflows;")
		pattern := &models.Pattern{
			Name: "lang_arutz_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang"},
				{Key: "archive_type", Value: "arutz"},
				{Key: "date"},
				{Key: "content_type"},
				{Key: "line"},
				{Key: "name"},
			},
			Extension: "mpg",
		}
		pattern.Save()

		workflow := &models.Workflow{
			PatternId:   sql.NullInt64{Int64: pattern.ID, Valid: true},
			EntryPoint:  "SiumAvoda",
			ContentType: sql.NullString{String: "*", Valid: true},
			Line:        sql.NullString{String: "*", Valid: true},
		}
		workflow.Save()

		file := &models.File{
			FileName:   "heb_arutz_2012-12-16_film_crossroads.mpg",
			TargetDir:  "targetDir",
			EntryPoint: "SiumAvoda",
			SourcePath: "path",
			//FullPath:   "../../test_files/heb_o_rav_achana_2015-10-13_lesson.mp4",
			//FullPath:   "../../test_files/heb_o_rav_2016-02-25_promo_congress_pre-roll.mpg",
			FullPath:   "../../test_files/mlt_o_rav_achana_2016-04-28_lesson.mpg",
			PatternId:  sql.NullInt64{Int64: pattern.ID, Valid: true},
			WorkflowId: sql.NullInt64{Int64: workflow.ID, Valid: true},
			Status:     models.HAS_WORKFLOW,
			Attributes: models.JSONB{"date": "2012-12-16", "lang": "heb", "line": "crossroads", "archive_type": "arutz", "content_type": "film"},
		}

		file.CreateVersion()

		Convey("When there is an existing validation", func() {
			workflow.Validations = models.StringSlice{"passedValidation", "failedValidation"}
			workflow.Save()
			Convey("It should run all workflow's validations as expected", func() {
				err := v.RunValidations(file)
				So(err, ShouldBeNil)
				So(file.ValidationResult, ShouldNotBeNil)
				/*
					 {
						 "failedValidation": {"Passed": false, "ErrorMessage": {}},
						 "passedValidation": {"Passed": true, "ErrorMessage": null}
					 }
				*/
				expectedValidationResult := models.JSONB{
					"passedValidation": models.ValidationResult{Passed: true, ErrorMessage: ""},
					"failedValidation": models.ValidationResult{Passed: false, ErrorMessage: "error"},
				}
				fmt.Printf("%v", file.ValidationResult["passedValidation"])

				So(file.ValidationResult["passedValidation"].(models.ValidationResult).Passed, ShouldEqual, expectedValidationResult["passedValidation"].(models.ValidationResult).Passed)
				So(file.ValidationResult["failedValidation"].(models.ValidationResult).Passed, ShouldEqual, expectedValidationResult["failedValidation"].(models.ValidationResult).Passed)
				So(file.ValidationResult["passedValidation"].(models.ValidationResult).ErrorMessage, ShouldEqual, expectedValidationResult["passedValidation"].(models.ValidationResult).ErrorMessage)
				So(file.ValidationResult["failedValidation"].(models.ValidationResult).ErrorMessage, ShouldEqual, expectedValidationResult["failedValidation"].(models.ValidationResult).ErrorMessage)
			})
			Convey("it should validate Exif properties correctly", func() {
				workflow.Validations = models.StringSlice{"checkExif"}
				workflow.Save()
				testConfig := []struct {
					fileName string
					passed   bool
					exif     models.Exif
				}{
					{"../../test_files/mlt_o_rav_achana_2016-04-28_lesson.mpg", true, models.Exif{
						FileType:          "MPEG",
						FileTypeExtension: "mpg",
						ImageWidth:        720,
						ImageHeight:       576,
						AspectRatio:       "16:9, 625 line, PAL",
						FrameRate:         "25 fps",
						VideoBitrate:      "12 Mbps",
						MPEGAudioVersion:  1,
						AudioLayer:        2,
						AudioBitrate:      "224 kbps",
						SampleRate:        48000,
						OriginalMedia:     true,
					}},
					{"../../test_files/heb_o_rav_2016-02-25_promo_congress_pre-roll.mpg", false, models.Exif{
						FileType:          "MPEG",
						FileTypeExtension: "mpg",
						ImageWidth:        730,
						ImageHeight:       576,
						AspectRatio:       "16:9, 625 line, PAL",
						FrameRate:         "25 fps",
						VideoBitrate:      "12 Mbps",
						MPEGAudioVersion:  1,
						AudioLayer:        2,
						AudioBitrate:      "224 kbps",
						SampleRate:        48000,
						OriginalMedia:     true,
					}},
				}

				for _, t := range testConfig {
					file := &models.File{
						FileName:   "heb_arutz_2012-12-16_film_crossroads.mpg",
						TargetDir:  "targetDir",
						EntryPoint: "SiumAvoda",
						SourcePath: "path",
						FullPath:   t.fileName,
						PatternId:  sql.NullInt64{Int64: pattern.ID, Valid: true},
						WorkflowId: sql.NullInt64{Int64: workflow.ID, Valid: true},
						Status:     models.HAS_WORKFLOW,
						Attributes: models.JSONB{"date": "2012-12-16", "lang": "heb", "line": "crossroads", "archive_type": "arutz", "content_type": "film"},
					}
					file.CreateVersion()
					workflow.Exif = t.exif
					workflow.Save()
					err := v.RunValidations(file)
					So(err, ShouldBeNil)
					So(file.ValidationResult, ShouldNotBeNil)
					So(file.ValidationResult["checkExif"].(models.ValidationResult).Passed, ShouldEqual, t.passed)
				}
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
