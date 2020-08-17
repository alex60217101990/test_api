package models

type Product struct {
	Base
	ChangeByUser uint        `sql:"type:integer;not null" json:"change_by_user"`
	Name         string      `sql:"type:varchar(100);not null" json:"name"`
	Popularity   uint32      `sql:"type:bigint" json:"popularity"`
	User         *User       `sql:"-" json:"user,omitempty"`
	Categories   []*Category `sql:"-" json:"categories,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

func (p Product) ConvertToQuery(s *SortedBy) (query string) {
	return convertToQuery(s, p)
}
