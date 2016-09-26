package models

import (
	"fmt"
	"regexp"
	"strings"
)

type Pattern struct {
	Model
	Name      string `sql:"not null;unique;size:255"`
	Priority  int
	Regexp    RegularX `sql:"type:varchar(255);not null;unique"`
	Extension string
	Workflows []Workflow
	Parts     Pairs `sql:"type:jsonb"`
}

type Patterns []Pattern

func (p *Pattern) FindOne() error {
	return db.First(p, p).Error
}

func (p *Pattern) Save() error {
	p.calculatePriorityField()

	if err := p.calculateRegexField(); err != nil {
		return err
	}
	return db.Save(p).Error
}

func (ps *Patterns) FindAll() error {
	return db.Find(ps).Error
}

func (ps *Patterns) FindAllByFileMatch(fileName string) error {
	return db.Where("? ~ regexp", fileName).Order("priority desc").Find(ps).Error
}

func (p *Pattern) calculatePriorityField() {
	for _, element := range p.Parts {
		if element.Value != "" {
			p.Priority += 1
		}
	}
}

func (p *Pattern) calculateRegexField() error {
	var parts []string

	for _, element := range p.Parts {
		patPart := &PatternPart{Key: element.Key}
		if err := patPart.FindOneByKey(); err != nil {
			return fmt.Errorf("PatternPart with key %q: %v", element.Key, err)
		}

		var value string
		if element.Value == "" {
			value = patPart.Value
		} else {
			value = element.Value
		}

		parts = append(parts, fmt.Sprintf("(%s)", value))
	}

	str := fmt.Sprintf("%s.(%s)", strings.Join(parts, "_"), p.Extension)
	p.Regexp.Regx, _ = regexp.Compile(str)

	return nil
}
