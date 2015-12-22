package workflow_manager

import (
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/models"
)

var AttachToWorkflow = file_manager.HandlerFunc(func(file *models.File) error {
	return nil
})
