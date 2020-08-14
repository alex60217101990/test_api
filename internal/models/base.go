package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	uuidEmptyStringValue = "00000000-0000-0000-0000-000000000000"
)

type Base struct {
	ID        int64      `sql:"primary_key type:int unsigned auto_increment" json:"id"`
	PublicID  uuid.UUID  `sql:"type:uuid;index" json:"public_id"`
	CreatedAt *time.Time `sql:"type:created_at;not null" json:"created_at"`
	UpdatedAt *time.Time `sql:"type:updated_at;not null" json:"updated_at"`
	DeletedAt *time.Time `sql:"type:deleted_at;index" json:"deleted_at"`
}

func (base *Base) BeforeCreate() (err error) {
	if base.PublicID.String() == uuidEmptyStringValue {
		base.PublicID, err = uuid.NewUUID()
	}
	return err
}
