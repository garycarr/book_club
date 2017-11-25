package util

import "github.com/garycarr/book_club/common"

// UtilIn is mainly for functions that need mocking
type UtilIn interface {
	CreateHashedPassword(string) (string, error)
	CheckHashedPassword(string, string) error
	CheckJSONToken(string) error
	CreateJSONToken(*common.User) (string, error)
}
