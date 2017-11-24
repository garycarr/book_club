package warehouse

import "github.com/garycarr/book_club/common"

// WarehouseIn ...
type WarehouseIn interface {
	CreateUser(common.RegisterRequest) (*common.User, error)
	GetUserWithEmail(string) (*common.User, error)
}
