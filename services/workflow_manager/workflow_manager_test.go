package workflow_manager_test

import (
	"log"
	"testing"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	wm "github.com/Bnei-Baruch/mms-file-manager/services/workflow_manager"
	"github.com/jinzhu/gorm"
	"github.com/Bnei-Baruch/mms-file-manager/test_helpers"
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	l  *log.Logger = logger.InitLogger(&logger.LogParams{LogMode: "screen", LogPrefix: "[WM-TEST] "})
	db *gorm.DB
)

func TestWorkflowSpec(t *testing.T) {
	db = test_helpers.SetupSpec()

	Convey("Setup", t, func() {

		SkipConvey("WorkflowManager", func() {
			Convey(" Describe Workflow saving", func() {
				Convey("When one pattern matched", func() {

					Convey("It should attach pattern to file", nil)
				})
			})
		})

		Convey("Descrive Workflow matching", func() {
			db.Exec("DELETE FROM patterns;")

			//TODO: content type/line type has variations - should test all

			Convey("When workflow is attached to file and passes validation", func() {
				Convey("It should attach workflow to file and mark file status as HAS_WORKFLOW", func() {
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
					err := workflow.Save()
					So(err, ShouldBeNil)

					file := &models.File{
						FileName:   "heb_arutz_2012-12-16_film_crossroads.mpg",
						TargetDir:  "targetDir",
						EntryPoint: "SiumAvoda",
						SourcePath: "path",
						PatternId:  sql.NullInt64{Int64: pattern.ID, Valid: true},
						Status:     models.HAS_PATTERN,
					}

					file.CreateVersion()
					err = wm.AttachToWorkflow(file)
					So(err, ShouldBeNil)
					So(file.WorkflowId.Int64, ShouldEqual, workflow.ID)
					So(file.Status, ShouldEqual, models.HAS_WORKFLOW)
				})
			})
			Convey("When workflow is attached to file but doesn't pass validation", func() {
				Convey("It should attach workflow to file and mark file status as HAS_NO_VALID_WORKFLOW", nil)
			})
			Convey("When workflow is not attached to file", func() {
				Convey("It should mark file status as HAS_NO_WORKFLOW", nil)
			})
		})
	})
}
