package models

type Product struct {
	Base
	ChangeByUser int         `sql:"type:integer;not null"`
	CategorieID  int         `sql:"type:integer;not null"`
	Name         string      `sql:"type:varchar(100);not null" json:"name"`
	Popularity   uint32      `sql:"type:bigint" json:"popularity"`
	User         *User       `sql:"-" json:"user,omitempty"`
	Categories   []*Category `sql:"-" json:"categories,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
