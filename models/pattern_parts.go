package models


type PatternPart struct {
	Model
	Key           string `sql:"not null;unique;size:255"`
	Value         string `sql:"size:255"`
}

type PatternParts []PatternPart


func (p *PatternPart) Save() error {
	return db.Save(p).Error
}

func (p *PatternPart) FindOne() error {
	return db.First(p, p).Error
}

func (p *PatternParts) FindAll() error {
	return db.Find(p).Error
}

