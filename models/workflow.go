package models


type Workflow struct {
	Model
	Pattern   Pattern
	PatternId int `sql:"index"`
	//	EntryPoint (label)
	//	MaterialType
	//	ArchiveType
	//	Language(s)
	//	ContentType
	//	HasName
	//	Line
	//	RequireCheck
}
