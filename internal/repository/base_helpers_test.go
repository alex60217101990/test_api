package repository

import (
	"fmt"
	"testing"

	"github.com/alex60217101990/test_api/internal/models"
	"github.com/google/uuid"
)

func TestConvertObjSliceToQueryStr(t *testing.T) {
	id1, err := uuid.NewUUID()
	if err != nil {
		t.Error(err)
		return
	}
	// id2, err := uuid.NewUUID()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	data := ConvertObjSliceToQueryStr(
		[]*models.Category{
			&models.Category{
				Base: models.Base{PublicID: id1},
			},
			// &models.Category{
			// 	Base: models.Base{PublicID: id2},
			// },
		},
	)
	if len(data) == 0 {
		t.Error("invalid type")
		return
	}
	fmt.Println(data)
}
