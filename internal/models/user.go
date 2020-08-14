package models

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	Base
	Username   string      `sql:"type:varchar(45);not null" json:"username"`
	Email      string      `sql:"type:varchar(100);not null" json:"email"`
	Password   string      `sql:"type:varchar(1000);not null" json:"password"`
	IsOnline   bool        `sql:"type:bool;not null" json:"is_online"`
	Categories []*Category `sql:"-" json:"categories"`
	Products   []*Product  `sql:"-" json:"products"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate() (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "before create user script")
		}
	}()

	if u.PublicID.String() == uuidEmptyStringValue {
		u.PublicID, err = uuid.NewUUID()
	}
	if err != nil {
		return err
	}

	// bts := make([]byte, 0)
	u.Password = hex.EncodeToString([]byte(u.Password))

	return err
}

func (u *User) AfterFind() (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "before after find user script")
		}
	}()

	if len(u.Password) > 0 {
		bts := make([]byte, 0)
		// bts, err = encrypt.DecryptWithPrivateKey(u.Password, configs.Conf.Keys.PubKeyRepo)
		bts, err = hex.DecodeString(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(bts)
	}

	return err
}
