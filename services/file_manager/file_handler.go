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
		EntryPoint: u.label,
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

var registrationHandler = HandlerFunc(func(file *models.File) (error) {
	file.Status = models.PENDING

	//First time the file should be created. All other places should use file.Update() instead!
	if err := file.CreateVersion(); err != nil {
		l.Println("Unable to create file:", file, " ERROR: ", err)
		return err
	}
	filePath := file.FilePath()
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		l.Println("Unable to create directory:", filePath, " ERROR: ", err)
		return err
	}
	file.FullPath = filepath.Join(filePath, file.FileName)
	os.Rename(file.SourcePath, file.FullPath)
	return nil
})

func (fm *FileManager) TargetDir() string {
	year, week := time.Now().ISOWeek()
	return filepath.Join(fm.TargetDirPrefix, strconv.Itoa(year), strconv.Itoa(week))
}
