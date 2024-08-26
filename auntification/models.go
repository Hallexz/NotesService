package auntification

import "errors"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (u *User) Validatee() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("неправильное имя или пароль")
	}
	return nil
}
