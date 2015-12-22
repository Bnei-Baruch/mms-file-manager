package workflow_manager

import (
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"database/sql"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
)

var AttachToPattern = file_manager.HandlerFunc(func(file *models.File) error {
	patterns := models.Patterns{}
	if err := patterns.FindAllByFileMatch(file.FileName); err != nil {
		return err
	}

	switch  len(patterns) {
	case 0:
		file.Status = models.NO_PATTERN
		file.PatternId = sql.NullInt64{}
	case 1:
		file.Status = models.HAS_PATTERN
		file.Pattern = patterns[0]
	default:
		if patterns[0].Priority == patterns[1].Priority {
			file.Status = models.MANY_PATTERNS
			file.PatternId = sql.NullInt64{}
		} else {
			file.Status = models.HAS_PATTERN
			file.Pattern = patterns[0]
		}
	}

	if file.Status == models.HAS_PATTERN {
		res := file.Pattern.Regexp.Regx.FindAllStringSubmatch(file.FileName, -1)[0][1:]
		attributes := make(models.JSONB, len(res))
		for index, el := range file.Pattern.Parts {
			attributes[el.Key] = res[index]
		}
		file.Attributes = attributes
	}

	return file.Save()
})