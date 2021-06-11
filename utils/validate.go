package utils

import (
	"github.com/badoux/checkmail"
)

func ValidateRegister(email, password string) bool {
	return validateEmail(email) && validatePassword(password)
}

func validateEmail(email string) bool {
	if email == "" {
		return false
	}

	var err error
	err = checkmail.ValidateFormat(email)
	if err != nil {
		return false
	}

	err = checkmail.ValidateHost(email)
	if err != nil {
		return false
	}

	return true
}

func validatePassword(password string) bool {
	if password == "" {
		return false
	}

	return true
}
