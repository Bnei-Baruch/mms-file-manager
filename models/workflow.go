package models
import "database/sql"


type Workflow struct {
	Model
	Pattern   Pattern
	PatternId  sql.NullInt64 `sql:"index"`
	EntryPoint  string
	ContentType string
	Line string
	//	MaterialType
	//	ArchiveType
	//	Language(s)
	//	HasName
	//	RequireCheck
}


func (w *Workflow) Save() error {
	return db.Save(w).Error
}

