package main

import "errors"

var errLoginUserAlreadyExists = errors.New("User already exists")
var errLoginUserNotFound = errors.New("User and password not found")
var errLoginUsernameNotPresent = errors.New("Username not present.")
var errLoginPasswordNotPresent = errors.New("Password not present.")
var errLoginUsernameAndPasswordNotPresent = errors.New("Username and password not present.")
