package models

type ContentType struct {
	Model
	Name string
}

func (c *ContentType) Save() error {
	return db.Save(c).Error
}

func (c *ContentType) Exists() (check bool, err error) {
	count := 0
	err = db.Model(ContentType{}).Where("name = ?", c.Name).Count(&count).Error
	check = count > 0
	return
}

func (c *ContentType) GetName() string {
	return c.Name
}