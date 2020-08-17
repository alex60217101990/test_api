package models

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (c *Credentials) BeforeQuery() (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "before create creeds script")
		}
	}()

	c.Password = hex.EncodeToString([]byte(c.Password))

	return err
}
