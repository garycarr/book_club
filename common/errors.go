package common

import "errors"

var ErrLoginEmailAndPasswordNotPresent = errors.New("Email and password not present")
var ErrLoginEmailNotPresent = errors.New("Email not present")
var ErrLoginPasswordNotPresent = errors.New("Password not present")
var ErrLoginUserAlreadyExists = errors.New("User already exists")
var ErrLoginUserNotFound = errors.New("Email and password not found or incorrect")

var ErrNewUserMissingFields = "Missing fields for new user:"
