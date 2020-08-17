package models

import (
	"reflect"
	"strings"
)

func convertToQuery(s *SortedBy, data interface{}) (query string) {
	defer func() {
		query = "ORDER BY " + query
		if query == "ORDER BY " {
			query = query + "created_at "
		}
		if s.Desc {
			query = query + " DESC "
		}
	}()

	v := reflect.ValueOf(data)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if len(strings.TrimSpace(s.FieldName)) > 0 &&
			strings.Split(t.Field(i).Tag.Get("json"), ",")[0] == s.FieldName {
			query = s.FieldName
			if s.Desc {
				query = query + " DESC "
			}
		}
	}

	return query
}
