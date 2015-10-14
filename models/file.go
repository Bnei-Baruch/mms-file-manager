package models

type File struct {
	Model
	TargetDir  string
	FileName  string
	Label string
	Status    string
//	Version string `sql:"index;type:varchar(100);unique" gorm:"column:kuku"`
	SourcePath string `sql:"-"` //will be ignored in DB
}