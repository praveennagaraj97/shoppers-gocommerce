package validators

import (
	"net/mail"
)

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}

	return nil
}
