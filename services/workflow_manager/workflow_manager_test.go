package workflow_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	wm "github.com/Bnei-Baruch/mms-file-manager/services/workflow_manager"
	"database/sql"
)

var _ = Describe("WorkflowManager", func() {
	PDescribe("Workflow saving", func() {
		Context("When one pattern matched", func() {

			It("should attach pattern to file", func() {

			})
		})
	})

	Describe("Workflow matching", func() {
		BeforeEach(func() {
			db.Exec("DELETE FROM patterns;")
		})

		//TODO: content type/line type has variations - should test all


		Context("When workflow is attached to file and passes validation", func() {
			It("should attach workflow to file and mark file status as HAS_WORKFLOW", func() {
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
				err := workflow.Save()
				Ω(err).ShouldNot(HaveOccurred())

				file := &models.File{
					FileName: "heb_arutz_2012-12-16_film_crossroads.mpg",
					TargetDir: "targetDir",
					EntryPoint: "SiumAvoda",
					SourcePath: "path",
					PatternId: sql.NullInt64{Int64: pattern.ID, Valid: true},
					Status: models.HAS_PATTERN,
				}

				err = wm.AttachToWorkflow(file)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(file.WorkflowId.Int64).Should(Equal(workflow.ID))
				Ω(file.Status).Should(Equal(models.HAS_WORKFLOW))

			})
		})
		Context("When workflow is attached to file but doesn't pass validation", func() {
			It("should attach workflow to file and mark file status as HAS_NO_VALID_WORKFLOW", func() { })
		})
		Context("When workflow is not attached to file", func() {
			It("should mark file status as HAS_NO_WORKFLOW", func() { })
		})
	})
})
