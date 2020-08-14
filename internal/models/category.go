package models

type Category struct {
	Base
	ChangeByUser int        `sql:"type:integer;not null"`
	Name         string     `sql:"type:varchar(100);not null" json:"name"`
	Popularity   int64      `sql:"type:bigint" json:"popularity"`
	User         *User      `sql:"-" json:"user,omitempty"`
	Products     []*Product `sql:"-" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "product_categories"
}
