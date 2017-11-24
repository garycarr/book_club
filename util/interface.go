package util

// UtilIn is an interface mainly for mocking
type UtilIn interface {
	GetCryptedPassword(string) (string, error)
}
