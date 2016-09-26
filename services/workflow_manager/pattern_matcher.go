package wmanager

import (
	"database/sql"

	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/services/file_manager"
)

//AttachToPattern function attaches patterns to file
var AttachToPattern = file_manager.HandlerFunc(func(file *models.File) (err error) {
	patterns := models.Patterns{}
	if err = patterns.FindAllByFileMatch(file.FileName); err != nil {
		return
	}

	defer func() {
		if err != nil {
			return
		}
		if err = file.Save(); err != nil {
			l.Printf("Problem updating file -  %s : %s", file.FileName, err)
		}
	}()

	switch len(patterns) {
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

	return
})
