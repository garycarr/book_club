package util

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

// CreateHashedPassword is mainly in a function so it can be mocked
func (u *Util) CreateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckHashedPassword ...
func (u *Util) CheckHashedPassword(dbPassword, givenPassword string) error {
	// Make sure the password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(givenPassword)); err != nil {
		return err
	}
	return nil
}
