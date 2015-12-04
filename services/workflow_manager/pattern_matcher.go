package workflow_manager
import (
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"database/sql"
)

func AttachToPattern(file *models.File) (bool, error) {
	patterns := models.Patterns{}
	if err := patterns.FindAllByFileMatch(file.FileName); err != nil {
		return false, err
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

	return file.Status == models.HAS_PATTERN, file.Save()
}