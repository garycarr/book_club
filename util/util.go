package util

import (
	"golang.org/x/crypto/bcrypt"
)

// Util ...
type Util struct{}

// NewUtil ...
func NewUtil() *Util {
	return &Util{}
}

const bcryptCost = 10

// GetCryptedPassword is mainly in a function so it can be mocked
func (u *Util) GetCryptedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
