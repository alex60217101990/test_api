package server

import (
	"fmt"

	"github.com/alex60217101990/test_api/internal/models"
)

func ValidateCreeds(creeds *models.Credentials) error {
	if creeds == nil {
		return fmt.Errorf("empty credentials data")
	}
	if len(creeds.Username) == 0 {
		return fmt.Errorf("invalid credentials data, empty 'username' field")
	}
	if len(creeds.Password) == 0 {
		return fmt.Errorf("invalid credentials data, empty 'password' field")
	}
	if len(creeds.Email) == 0 {
		return fmt.Errorf("invalid credentials data, empty 'email' field")
	}
	return nil
}
