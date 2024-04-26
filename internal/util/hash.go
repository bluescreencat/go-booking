package util

import (
	"golang.org/x/crypto/bcrypt"
)

func (u *utility) HashPassword(password string, round int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), round)
	return string(hash), err
}

func (u *utility) ComparePassword(hashPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err
}
