package repository

import (
	"fmt"
	"reflect"
	"strings"
)

func ConvertObjSliceToQueryStr(slice interface{}) (str string) {
	s := reflect.ValueOf(slice)
	delimiter := ", "
	if s.Kind() != reflect.Slice {
		return str
	}

	for i := 0; i < s.Len(); i++ {
		if base, ok := s.Index(i).Interface().(HasRelations); ok {
			str = fmt.Sprintf(`%s'%s'%s`, str, base.GetPublicID(), delimiter)
		}
	}

	return strings.Trim(str, delimiter)
}
