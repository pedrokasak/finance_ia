package auth

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmailFormat(email string) error {
		if !emailRegex.MatchString(email) {
				return errors.New("formato de email inválido")
		}
		return nil
}

func ValidateLoginFields(email, password string) error {
		if email == "" || password == "" {
				return errors.New("email e senha são obrigatórios")
		}
		if err := ValidateEmailFormat(email); err != nil {
				return err
		}
		return nil
}

func ErrIncorrectPassword() error {
    return errors.New("senha incorreta")
}
