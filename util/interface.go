package util

// UtilIn is mainly for functions that need mocking
type UtilIn interface {
	CreateHashedPassword(string) (string, error)
	CheckHashedPassword(string, string) error
}
