package util

import (
	"github.com/stretchr/testify/mock"
)

// MockUtil implements the Warehouse interface for the purpose of testing
type MockUtil struct {
	mock.Mock
}

// GetCryptedPassword is used to assert the method is called
func (mw *MockUtil) GetCryptedPassword(password string) (string, error) {
	args := mw.Called(password)
	return args.Get(0).(string), args.Error(1)
}
