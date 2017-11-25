package util

import (
	"github.com/garycarr/book_club/common"
	"github.com/stretchr/testify/mock"
)

// MockUtil implements the Warehouse interface for the purpose of testing
type MockUtil struct {
	mock.Mock
}

// CreateHashedPassword is used to assert the method is called
func (mw *MockUtil) CreateHashedPassword(password string) (string, error) {
	args := mw.Called(password)
	return args.Get(0).(string), args.Error(1)
}

// CheckHashedPassword is used to assert the method is called
func (mw *MockUtil) CheckHashedPassword(dbPassword, givenPassword string) error {
	args := mw.Called(dbPassword, givenPassword)
	return args.Error(0)
}

// CreateJSONToken is used to assert the method is called
func (mw *MockUtil) CreateJSONToken(user *common.User) (string, error) {
	args := mw.Called(user)
	return args.Get(0).(string), args.Error(1)
}

// CheckJSONToken is used to assert the method is called
func (mw *MockUtil) CheckJSONToken(token string) error {
	args := mw.Called(token)
	return args.Error(0)
}
