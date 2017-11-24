package common

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	JWTSecret     = "Environmentalise"
	JWTIssuer     = "Me"
	JWTExpiration = time.Duration(1 * time.Hour)
)

// CustomJWTClaims ..
type CustomJWTClaims struct {
	DisplayName string `json:"displayName"`
	jwt.StandardClaims
}

// CreateJSONToken ..
func CreateJSONToken(u *User) (string, error) {
	// Create the JSON token as the login is valid
	claims := &CustomJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTExpiration).Unix(),
			Issuer:    JWTIssuer,
			IssuedAt:  time.Now().Unix(),
			Id:        u.ID,
		},
		DisplayName: u.DisplayName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
