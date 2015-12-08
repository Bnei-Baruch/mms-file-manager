package workflow_manager_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("WorkflowManager", func() {
	PDescribe("Workflow saving", func() {
		Context("When one pattern matched", func() {

			It("should attach pattern to file", func() {

			})
		})
	})

	PDescribe("Workflow matching", func() {
		Context("When workflow is attached to file and passes validation", func() {
			It("should attach workflow to file", func() { })
			It("should mark file status as HAS_WORKFLOW", func() { })
		})
		Context("When workflow is attached to file but doesn't pass validation", func() {
			It("should attach workflow to file", func() { })
			It("should mark file status as NOT_VALID", func() { })
		})
		Context("When workflow is not attached to file", func() {
			It("should mark file status as HAS_NO_WORKFLOW", func() { })
		})
	})
})
