package user

import (
	"errors"
	"fmt"

	"github.com/badoux/checkmail"
)

func ValidateUserFields(user *User) error {
	fmt.Print("user: ", user)
		if user.FirstName == "" || user.Email == "" {
				return errors.New("data user not found")
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

func ValidateEmailFormat(email string) error {
		if email == "" {
				return errors.New("email é obrigatório")
		}
		err := checkmail.ValidateFormat(email)
		if err != nil {
				return errors.New("formato de email inválido")
		}
		return nil
}