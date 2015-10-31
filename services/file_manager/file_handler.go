package file_manager
import (
	"path/filepath"
	"os"
	"time"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"github.com/Bnei-Baruch/mms-file-manager/utils"
	"strconv"
)

func (fm *FileManager) handler(u updateMsg) {
	targetDir := fm.TargetDir()
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		l.Println("####################### Unable to create directory:", targetDir, " ERROR: ", err)
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
	file.Status = models.PENDING
	if err := file.Save(); err != nil {
		l.Println("Unable to save file:", file, " ERROR: ", err)
		return err
	}
	filePath := file.FilePath()
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		l.Println("Unable to create directory:", filePath, " ERROR: ", err)
		return err
	}

	os.Rename(file.SourcePath, filepath.Join(filePath, file.FileName))
	return nil
}

func (fm *FileManager) TargetDir() string {
	year, week := time.Now().ISOWeek()
	return filepath.Join(fm.TargetDirPrefix, strconv.Itoa(year), strconv.Itoa(week))
}
