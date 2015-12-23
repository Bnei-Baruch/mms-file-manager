package workflow_manager

import (
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
	"log"
)

var l *log.Logger = nil

func init() {
	l = logger.InitLogger(&logger.LogParams{LogPrefix: "[WM] "})
}

var AttachToWorkflow = file_manager.HandlerFunc(func(file *models.File) (err error) {

	defer func() {
		if err = file.Save(); err != nil {
			l.Printf("Problem saving file -  %s : %s", file.FileName, err)
		}
	}()

	workflows := models.Workflows{}
	if err = workflows.FindAllByEntryPointAndPatternId(file.PatternId, file.EntryPoint); err != nil {
		l.Printf("Problem attaching workflow to file -  %s : %s", file.FileName, err)
		file.Status = models.HAS_NO_WORKFLOW
		return
	}

	switch len(workflows) {
	case 0:
		file.Status = models.HAS_NO_WORKFLOW
	case 1:
		file.Status = models.HAS_WORKFLOW
		file.Workflow = workflows[0]
		//TODO: continue check content_type + line
	default:

		if contentType, ok := file.Attributes["content_type"]; ok {
		}
	}

	return
})