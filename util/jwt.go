package util

import (
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garycarr/book_club/common"
)

const (
	JWTSecret     = "Environmentalise"
	jwtIssuer     = "Me"
	jwtExpiration = time.Duration(1 * time.Hour)
)

// customJWTClaims ..
type customJWTClaims struct {
	DisplayName string `json:"displayName"`
	jwt.StandardClaims
}

// CreateJSONToken ..
func (u *Util) CreateJSONToken(user *common.User) (string, error) {
	// Create the JSON token as the login is valid
	claims := &customJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
			Issuer:    jwtIssuer,
			IssuedAt:  time.Now().Unix(),
			Id:        user.ID,
		},
		DisplayName: user.DisplayName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CheckJSONToken ...
func (u *Util) CheckJSONToken(token string) error {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return common.ErrJSONTokenNoBearer
	}
	jwToken := strings.TrimPrefix(token, "Bearer ")
	_, err := jwt.Parse(jwToken, func(jwToken *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}
