package validator_test

import (
	"database/sql"
	"testing"

	"log"
	"os"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"github.com/Bnei-Baruch/mms-file-manager/services/validator"
	"github.com/Bnei-Baruch/mms-file-manager/test_helpers"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/gomega"
)

var (
	l  *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[VAL-TEST] "})
	db *gorm.DB
)

func TestMain(m *testing.M) {

	db = test_helpers.SetupSpec()

	//l.Println("TestValidators start time", time.Now())
	code := m.Run()
	//l.Println("TestValidators end time", time.Now())
	os.Exit(code)
}

//TestValidators tests Workflow with list of validations
func TestValidators(t *testing.T) {
	RegisterTestingT(t)
	t.Run("test#1 When there is an existing validation", func(t *testing.T) {
		l.Println("kuku2222222")
		t.Run("test#1.1 It should run all workflow's validations as expected", func(t *testing.T) {
			o := &testObject{}
			o.prepareWorkflow()
			err := validator.RunValidations(o.file)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(o.file.ValidationResult).ShouldNot(BeEmpty())
			expectedValidationResult := models.JSONB{
				"passedValidation": models.ValidationResult{Passed: true, ErrorMessage: ""},
				"failedValidation": models.ValidationResult{Passed: false, ErrorMessage: "error"},
			}
			Ω(expectedValidationResult["passedValidation"].(models.ValidationResult).Passed).Should(BeEquivalentTo(o.file.ValidationResult["passedValidation"].(models.ValidationResult).Passed))
			Ω(expectedValidationResult["failedValidation"].(models.ValidationResult).Passed).Should(BeEquivalentTo(o.file.ValidationResult["failedValidation"].(models.ValidationResult).Passed))

			Ω(expectedValidationResult["passedValidation"].(models.ValidationResult).ErrorMessage).Should(BeEquivalentTo(o.file.ValidationResult["passedValidation"].(models.ValidationResult).ErrorMessage))
			Ω(expectedValidationResult["failedValidation"].(models.ValidationResult).ErrorMessage).Should(BeEquivalentTo(o.file.ValidationResult["failedValidation"].(models.ValidationResult).ErrorMessage))
		})
		t.Run("test#1.2 it should validate Exif properties correctly", func(t *testing.T) {
			o := &testObject{}
			o.prepareWorkflow()
			o.workflow.Validations = models.StringSlice{"checkExif"}
			o.workflow.Save()
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
					PatternId:  sql.NullInt64{Int64: o.pattern.ID, Valid: true},
					WorkflowId: sql.NullInt64{Int64: o.workflow.ID, Valid: true},
					Status:     models.HAS_WORKFLOW,
					Attributes: models.JSONB{"date": "2012-12-16", "lang": "heb", "line": "crossroads", "archive_type": "arutz", "content_type": "film"},
				}
				file.CreateVersion()
				o.workflow.Exif = t.exif
				o.workflow.Save()
				err := validator.RunValidations(file)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(file.ValidationResult).ShouldNot(BeNil())

				Ω(file.ValidationResult["checkExif"].(models.ValidationResult).Passed).Should(BeEquivalentTo(t.passed))
			}

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
}

type testObject struct {
	file     *models.File
	pattern  *models.Pattern
	workflow *models.Workflow
}

func (o *testObject) prepareWorkflow() {
	db.Exec("DELETE FROM patterns; DELETE FROM files; DELETE FROM workflows;")
	o.pattern = &models.Pattern{
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
	o.pattern.Save()

	o.workflow = &models.Workflow{
		PatternId:   sql.NullInt64{Int64: o.pattern.ID, Valid: true},
		EntryPoint:  "SiumAvoda",
		ContentType: sql.NullString{String: "*", Valid: true},
		Line:        sql.NullString{String: "*", Valid: true},
	}
	o.workflow.Save()
	o.file = &models.File{
		FileName:   "mlt_o_rav_achana_2016-04-28_lesson.mpg",
		TargetDir:  "targetDir",
		EntryPoint: "SiumAvoda",
		SourcePath: "path",
		//FullPath:   "../../test_files/heb_o_rav_achana_2015-10-13_lesson.mp4",
		//FullPath:   "../../test_files/heb_o_rav_2016-02-25_promo_congress_pre-roll.mpg",
		FullPath:   "../../test_files/mlt_o_rav_achana_2016-04-28_lesson.mpg",
		PatternId:  sql.NullInt64{Int64: o.pattern.ID, Valid: true},
		WorkflowId: sql.NullInt64{Int64: o.workflow.ID, Valid: true},
		Status:     models.HAS_WORKFLOW,
		Attributes: models.JSONB{"date": "2012-12-16", "lang": "heb", "line": "crossroads", "archive_type": "arutz", "content_type": "film"},
	}

	o.file.CreateVersion()
	o.workflow.Validations = models.StringSlice{"passedValidation", "failedValidation"}
	o.workflow.Save()
	l.Println("kuku1111111")
}
