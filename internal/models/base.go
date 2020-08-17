package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	uuidEmptyStringValue = "00000000-0000-0000-0000-000000000000"
)

type Base struct {
	ID        int64      `sql:"primary_key type:int unsigned auto_increment" json:"-"`
	PublicID  uuid.UUID  `sql:"type:uuid;index" json:"public_id"`
	CreatedAt *time.Time `sql:"created_at" time_format:"sql_date" time_utc:"true" json:"created_at"`
	UpdatedAt *time.Time `sql:"updated_at" time_format:"sql_date" time_utc:"true" json:"updated_at"`
	DeletedAt *time.Time `sql:"deleted_at" time_format:"sql_date" time_utc:"true" json:"deleted_at"`
}

func (base Base) GetPublicID() string {
	return base.PublicID.String()
}

func (base *Base) FromStr(str string) (err error) {
	base.PublicID, err = uuid.Parse(str)
	return err
}

func (base *Base) IsEmpty() bool {
	return base.PublicID.String() == uuidEmptyStringValue
}

func (base *Base) BeforeCreate() (err error) {
	if base.PublicID.String() == uuidEmptyStringValue {
		base.PublicID, err = uuid.NewUUID()
	}
	return err
}
