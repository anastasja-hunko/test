package model

import "golang.org/x/crypto/bcrypt"

//user struct
type User struct {
	Login     string `bson:"login,"`
	Password  string `bson:"password,"`
	Documents []interface{}
}

//hash user's password before create to db
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

//compare password while autorizing
func (u *User) ComparePasswords(password string) bool {
	return checkPasswordHash(u.Password, password)
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

//hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
