package models

import (
	"reflect"
	"strings"
)

type Category struct {
	Base
	ChangeByUser uint       `sql:"type:integer;not null" json:"change_by_user"`
	Name         string     `sql:"type:varchar(100);not null" json:"name"`
	Popularity   int64      `sql:"type:bigint" json:"popularity"`
	User         *User      `sql:"-" json:"user,omitempty"`
	Products     []*Product `sql:"-" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "product_categories"
}

func (c Category) ConvertToQuery(s *SortedBy) (query string) {
	defer func() {
		query = "ORDER BY " + query
		if len(query) == 0 {
			query = query + "created_at"
		}
	}()

	v := reflect.ValueOf(c)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if strings.Split(t.Field(i).Tag.Get("json"), ",")[0] == s.FieldName {
			query = s.FieldName
			if s.Desc {
				query = query + " DESC "
			}
		}
	}

	return query
}
