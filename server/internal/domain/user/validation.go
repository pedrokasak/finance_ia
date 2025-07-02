package user

import "errors"

func ValidateCreateUserFields(firstName, lastName, email, password string) error {
    if firstName == "" || lastName == "" || email == "" || password == "" {
        return errors.New("dados obrigatórios ausentes")
    }
    return nil
}

func ValidateUpdateUserFields(user *User) error {
		if user.ID == ([16]byte{}) {
				return errors.New("ID inválido")
		}
		if user.FirstName == "" || user.Email == "" {
				return errors.New("dados obrigatórios ausentes")
		}
		return nil
}

func ValidateLoginFields(email, password string) error {
		if email == "" || password == "" {
				return errors.New("email e senha são obrigatórios")
		}
		return nil
}