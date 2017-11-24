package main

import "errors"

var errLoginEmailAndPasswordNotPresent = errors.New("Email and password not present")
var errLoginEmailNotPresent = errors.New("Email not present")
var errLoginPasswordNotPresent = errors.New("Password not present")
var errLoginUserAlreadyExists = errors.New("User already exists")
var errLoginUserNotFound = errors.New("Email and password not found or incorrect")

var errNewUserMissingFields = "Missing fields for new user:"
