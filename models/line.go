package models

type Line struct {
	Model
	Name string
}

func (l *Line) Save() error {
	return db.Save(l).Error
}

func (l *Line) Exists() (check bool, err error) {
	count := 0
	err = db.Model(Line{}).Where("name = ?", l.Name).Count(&count).Error
	check = count > 0
	return
}

func (l *Line) GetName() string {
	return l.Name
}
