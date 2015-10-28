package models
import "fmt"


type File struct {
	Model
	TargetDir  string
	FileName   string
	Label      string           // incoming source of file (i.e. ingest, etc.)
	Status     string
	Version    int64
								//	Version string `sql:"index;type:varchar(100);unique" gorm:"column:kuku"`
	SourcePath string `sql:"-"` //will be ignored in DB
}

func (f *File) Load() error {
	return db.Where(f).First(f).Error
}

func (f *File) Save() error {
	// check if present in DB
	//	if not just add file record and rename + move file
	//	if present
	//	1. if file exist create second version
	//	2. if file doesn't exist - rename + move file
	var lastFile File
	if db.Where(File{FileName: f.FileName}).Order("version desc").First(&lastFile).RecordNotFound() {
		f.Version = 1
	} else {
		f.Version = lastFile.Version + 1
	}
	return db.Create(f).Error
}

func (f *File) FilePath() string {
	return fmt.Sprintf("%s/v%02d", f.TargetDir, f.Version)
}
