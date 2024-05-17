package model

import (
	"errors"
)

type User struct {
	login    string
	password string
}

func NewUser(login, password string) (*User, error) {
	user := User{
		login:    login,
		password: password,
	}

	if !user.Validate() {
		return nil, errors.New("не верный логин или пароль")
	}

	return &user, nil
}

func (u User) Validate() bool {
	if len(u.login) == 0 || len(u.password) == 0 {
		return false
	}

	return true
}
