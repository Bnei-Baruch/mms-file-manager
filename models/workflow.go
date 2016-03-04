package models

import (
	"database/sql"
	"reflect"
	"github.com/sthorne/reflections"
)

type Workflow struct {
	Model
	Pattern     Pattern
	PatternId   sql.NullInt64 `sql:"index"`
	EntryPoint  string
	ContentType sql.NullString
	Line        sql.NullString
	Validations StringSlice `sql:"type:varchar[]"`
	//	MaterialType
	//	ArchiveType
	//	Language(s)
	//	HasName
	//	RequireCheck
}
type LookupField interface {
	GetName() string
	Exists() (bool, error)
}

type Workflows []Workflow

func (w *Workflow) Save() error {
	return db.Save(w).Error
}

func (ws *Workflows) detectWorkflowsBy(fieldName string, checkInterface LookupField) (err error) {
	result := make(map[int]Workflows, 3)
	checkValue := checkInterface.GetName()
	var (
		check bool
		workflowValue interface{}
	)

	for _, w := range *ws {
		workflowValue, err = reflections.GetFieldTag(w, fieldName, "String")
		if err != nil {
			ws = nil
			return
		}

		switch (workflowValue) {
		case "": // Priority 3;
			result[3] = append(result[3], w)
		case "*": // Priority 2
			if check, err = checkInterface.Exists(); err != nil {
				l.Printf("Could not look for: %q, in content_types. %s\n", checkValue, err)
				ws = nil
				return
			}
			if check {
				result[2] = append(result[2], w)
			}
		default: // Priority 1
			if workflowValue == checkValue {
				result[1] = append(result[1], w)
			}
		}
	}

	for _, i := range []int{1, 2, 3} {
		if ok, _ := result[i]; ok != nil {
			r := result[i]
			ws = &r
			return
		}
	}
	ws = nil
	return
}

func (ws *Workflows) FindAllByEntryPointAndPatternId(patternId sql.NullInt64, entryPoint string) error {
	return db.Where("pattern_id = ? and entry_point = ?", patternId.Int64, entryPoint).Find(ws).Error
}

func (ws *Workflows) DetectWorkflows(contentType, line string) (err error) {
	if err = ws.detectWorkflowsBy("ContentType", &ContentType{Name: contentType}); err != nil {
		return
	}

	if ws == nil {
		return
	}

	err = ws.detectWorkflowsBy("Line", &Line{Name: line})

	return
}

func (w *Workflow) getStringField(field string) string {
	r := reflect.ValueOf(w)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}