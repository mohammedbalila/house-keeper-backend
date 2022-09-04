package db

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

type User struct {
	Base
	Id            string `json:"id"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	FirebaseToken string `json:"firebaseToken"`
	Password      string `json:"-"`

	tableName struct{} `pg:"api.user"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%s %s>", u.Id, u.FullName)
}

func (u *User) HashPassword() error {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(u.Password))
	if err != nil {
		return err
	}
	u.Password = string(encoded)
	return nil
}

func (u User) VerifyPassword(password string) (bool, error) {
	match, err := argon2.VerifyEncoded([]byte(password), []byte(u.Password))
	if err != nil {
		return false, err
	}
	if !match {
		return false, nil
	}
	return true, nil
}
