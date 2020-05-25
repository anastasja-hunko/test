package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	Login     string `bson:"login,"`
	Password  string `bson:"password,"`
	Documents []interface{}
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := hashPassword(u.Password)

		if err != nil {
			return err
		}

		u.Password = enc
	}

	return nil
}

func (u *User) ComparePasswords(password string) bool {
	return checkPasswordHash(u.Password, password)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
