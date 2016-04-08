package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
)

type Status string

const (
	PENDING Status = "PENDING"
	NEW     Status = "NEW"

	HAS_PATTERN   Status = "HAS_PATTERN"
	NO_PATTERN    Status = "NO_PATTERN"
	MANY_PATTERNS Status = "MANY_PATTERNS"

	HAS_WORKFLOW          Status = "HAS_WORKFLOW"
	HAS_NO_WORKFLOW       Status = "HAS_NO_WORKFLOW"
	HAS_MANY_WORKFLOWS    Status = "HAS_MANY_WORKFLOWS"
	HAS_NO_VALID_WORKFLOW Status = "HAS_NO_VALID_WORKFLOW"
)

type ValidationResult struct {
	Passed       bool
	ErrorMessage string
}

//type ValidationResults map[string]ValidationResult

type File struct {
	Model
	TargetDir  string
	FileName   string
	FullPath   string `sql:"type:text"` // - TargetDir/Version/FileName
	EntryPoint string // incoming source of file (i.e. ingest, etc.)
	Status     Status `sql:"type:varchar(30)"`
	Error      string `sql:"type:text"`
	Version    int64
	//	Version string `sql:"index;type:varchar(100);unique" gorm:"column:kuku"`
	SourcePath       string `sql:"-"` //will be ignored in DB
	Pattern          Pattern
	PatternId        sql.NullInt64 `sql:"index"`
	Attributes       JSONB         `sql:"type:jsonb"` // parsed attributes out of file name
	Workflow         Workflow
	WorkflowId       sql.NullInt64 `sql:"index"`
	ValidationResult JSONB         `sql:"type:jsonb"` // map[string] struct{passed bool, err_message string}
}

func (f *File) Load() (err error) {
	err = db.Where(f).Preload("Workflow").Preload("Pattern").First(f).Error
	if f.ValidationResult == nil {
		f.ValidationResult = make(JSONB)
	}

	return err
}

func (f *File) Save() error {
	return db.Save(f).Error
}

func (f *File) CreateVersion() error {
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

func (u *Status) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*u = Status(asBytes)
	return nil
}

func (u Status) Value() (driver.Value, error) {
	return string(u), nil
}
