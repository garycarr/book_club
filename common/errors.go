package common

import "errors"

var (
	ErrLoginEmailAndPasswordNotPresent = errors.New("Email and password not present")
	ErrLoginEmailNotPresent            = errors.New("Email not present")
	ErrLoginPasswordNotPresent         = errors.New("Password not present")
	ErrLoginUserAlreadyExists          = errors.New("User already exists")
	ErrLoginUserNotFound               = errors.New("Email and password not found or incorrect")

	ErrNewUserMissingFields = "Missing fields for new user:"
)
