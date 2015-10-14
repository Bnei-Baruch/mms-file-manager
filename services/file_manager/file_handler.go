package file_manager
import (
	"path/filepath"
	"os"
	"time"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/utils"
)

func (fm *FileManager) handler(u updateMsg) {
	targetDir := fm.TargetDir()
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return
	}

	file := &models.File{
		FileName:    filepath.Base(u.filePath),
		TargetDir: targetDir,
		Label: u.label,
		SourcePath: u.filePath,
	}

	for _, handler := range fm.handlers {
		if err := handler(file); err != nil {
			handlerName := utils.GetFunctionName(handler)
			l.Printf("Failed to call handler %q for file: %#v \n\tError: %v\n", handlerName, file, err)
			return
		}
	}
}

func registrationHandler(file *models.File) (error) {
// check if present in DB
//	if not just add file record and rename + move file
//	if present
//	1. if file exist create second version
//	2. if file doesn't exist - rename + move file
	fileName := file.FileName
	os.Rename(file.SourcePath, filepath.Join(file.TargetDir, fileName))
	return nil
}

func (fm *FileManager) TargetDir() string {
	year, week := time.Now().ISOWeek()
	return filepath.Join(fm.TargetDirPrefix, string(year), string(week))
}
