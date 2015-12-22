package models
import (
	"time"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"database/sql/driver"
	"regexp"
)

type Model struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//	DeletedAt *time.Time `sql:"index"`
}

var db *gorm.DB

func New(dbConn *gorm.DB) {
	db = dbConn
}

//Support for string slices stored as `sql: "text"`

type Strings []string

func (l Strings) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l *Strings) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), l)
}


type RegularX struct {
	Regx *regexp.Regexp
}

// Save regular expression to DB as string
func (r RegularX) Value() (driver.Value, error) {
	return r.Regx.String(), nil
}

// Read regular expression string from DB and compile it
func (r *RegularX) Scan(input interface{}) (err error) {
	str := string(input.([]byte))
	r.Regx, err = regexp.Compile(str)
	return
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}


type Pairs []struct {
	Key   string `json:"key"`
	Value string `json:"value;omitempty"`
}

func (j Pairs) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *Pairs) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}