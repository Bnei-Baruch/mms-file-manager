package models

import (
	"database/sql"
	"reflect"
)

type Workflow struct {
	Model
	Priority    int `sql:"-"`
	Pattern     Pattern
	PatternId   sql.NullInt64 `sql:"index"`
	EntryPoint  string
	ContentType sql.NullString
	Line        string
	//	MaterialType
	//	ArchiveType
	//	Language(s)
	//	HasName
	//	RequireCheck
}

type Workflows []Workflow

func (w *Workflow) Save() error {
	return db.Save(w).Error
}

func (w *Workflow) detectWorkflowByContentType(fieldName, fieldValue string, lookupInterface interface{}) (check bool) {
	var err error
	lookupStruct := reflect.TypeOf(lookupInterface).
	value := w.getStringField(fieldName)
	switch value {
	case "": // Priority 3
		check = true
		w.Priority += 3
	case "*": // Priority 2
		check, err = lookupInterface.(lookupStruct).Exists(fieldValue)
		if err != nil {
			l.Printf("Could not look for: %q, in content_types. %s\n", fieldValue, err)
			check = false
		}
		if check {
			w.Priority += 2
		}
	default: // Priority 1
		check = value == fieldValue
		if check {
			w.Priority += 1
		}
	}

	return
}

func (ws *Workflows) FindAllByEntryPointAndPatternId(patternId sql.NullInt64, entryPoint string) error {
	//TODO: if workflows should have priority we need to order by + TEST
	return db.Where("pattern_id = ? and entry_point = ?", patternId, entryPoint).Find(ws).Error
}

func (ws *Workflows) DetectWorkflows(contentType, line string) {
	result := &Workflows{}
	for _, w := range ws {
		w.Priority = 0
		if w.detectWorkflowByContentType(contentType) {
			append(result, w)
		}
	}
	ws = result
}

func (w *Workflow) getStringField(field string) string {
	r := reflect.ValueOf(w)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}