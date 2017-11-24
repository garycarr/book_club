package common

import (
	"fmt"
	"strings"
)

const bcryptCost = 10

// LoginRequest is the data needed to make a login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest is the information needed to register a new user
type RegisterRequest struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

// User ...
type User struct {
	Email       string
	ID          string
	DisplayName string
	Password    string
}

// ValidateRequest ..
func (lr LoginRequest) ValidateRequest() error {
	if lr.Email == "" && lr.Password == "" {
		return ErrLoginEmailAndPasswordNotPresent
	}
	if lr.Email == "" {
		return ErrLoginEmailNotPresent
	}
	if lr.Password == "" {
		return ErrLoginPasswordNotPresent
	}
	return nil
}

// ValidateNewUserRequest ...
func (rr RegisterRequest) ValidateNewUserRequest() error {
	var missingFields string
	if rr.DisplayName == "" {
		missingFields = "displayName,"
	}
	if rr.Password == "" {
		if missingFields != "" {
			missingFields = fmt.Sprintf("%s ", missingFields)
		}
		missingFields += "password,"
	}
	if rr.Email == "" {
		if missingFields != "" {
			missingFields = fmt.Sprintf("%s ", missingFields)
		}
		missingFields += "email,"
	}
	if missingFields != "" {
		missingFields = strings.TrimRight(missingFields, ",")
		return fmt.Errorf(fmt.Sprintf("%s %s", ErrNewUserMissingFields, missingFields))
	}
	return nil
}
