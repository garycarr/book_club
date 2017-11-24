package warehouse

import (
	"github.com/garycarr/book_club/common"
	"github.com/stretchr/testify/mock"
)

// MockWarehouse implements the Warehouse interface for the purpose of testing
type MockWarehouse struct {
	mock.Mock
}

// CreateUser is used to assert the method is called
func (mw *MockWarehouse) CreateUser(rr common.RegisterRequest) (*common.User, error) {
	args := mw.Called(rr)
	return args.Get(0).(*common.User), args.Error(1)
}

// GetUserWithEmail is used to assert the method is called
func (mw *MockWarehouse) GetUserWithEmail(email string) (*common.User, error) {
	args := mw.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.User), args.Error(1)
}

// Close is used to assert the method is called
func (mw *MockWarehouse) Close() {}
