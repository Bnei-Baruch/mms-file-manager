package wmanager

import (
	"fmt"
	"log"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
	"github.com/Bnei-Baruch/mms-file-manager/services/logger"
)

var l *log.Logger

func init() {
	l = logger.InitLogger(&logger.LogParams{LogPrefix: "[WM] "})
}

//AttachToWorkflow attaches workflow to a file
var AttachToWorkflow = file_manager.HandlerFunc(func(file *models.File) (err error) {

	defer func() {
		if err != nil {
			return
		}
		if err = file.Save(); err != nil {
			l.Printf("Problem updating file -  %s : %s", file.FileName, err)
		}
	}()

	workflows := models.Workflows{}
	if err = workflows.FindAllByEntryPointAndPatternId(file.PatternId, file.EntryPoint); err != nil {
		file.Status = models.HAS_NO_WORKFLOW
		file.Error = fmt.Sprintf("Problem attaching workflow to file -  %s : %s", file.FileName, err)
		l.Println(file.Error)
		return
	}

	if len(workflows) == 0 {
		file.Status = models.HAS_NO_WORKFLOW
		file.Error = fmt.Sprintf("No workflows matched by entry point %q and pattern %v", file.EntryPoint, file.PatternId)
		l.Println(file.Error)
		return
	}

	detectWorkflows(file, workflows)
	return
})

func detectWorkflows(file *models.File, workflows models.Workflows) {
	contentType, _ := file.Attributes["content_type"].(string)
	line, _ := file.Attributes["line"].(string)
	if err := workflows.DetectWorkflows(contentType, line); err != nil {
		file.Status = models.HAS_NO_WORKFLOW
		file.Error = fmt.Sprintf("Could not detect workflows for file: %#v. %s\n", file, err)
		l.Println(file.Error)
	}
	if len(workflows) > 0 {
		if len(workflows) == 1 {
			file.Status = models.HAS_WORKFLOW
			file.Workflow = workflows[0]

		} else {
			file.Status = models.HAS_MANY_WORKFLOWS
			file.Error = fmt.Sprintf("Matched more than one workflow\n: %#v", workflows)
			l.Println(file.Error)
		}
	} else {
		file.Status = models.HAS_NO_WORKFLOW
		file.Error = fmt.Sprintf("No workflows matched by content_type %q and line %q", contentType, line)
		l.Println(file.Error)
	}
}
