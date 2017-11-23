package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func setupTest(req *http.Request) (*app, *httptest.ResponseRecorder) {
	rr := httptest.NewRecorder()
	a := app{}
	a.initialize("test_config.json")
	return &a, rr
}

func checkJWT(t *testing.T, expectedClaims customJWTClaims, tokenString string, testDescription string) error {
	claims := customJWTClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return fmt.Errorf("Error parsing jwt for test %q: %v", testDescription, err)
	}

	assert.Equal(t, expectedClaims.Id, claims.Id, testDescription)
	assert.Equal(t, expectedClaims.Issuer, claims.Issuer, testDescription)
	assert.Equal(t, expectedClaims.Username, claims.Username, testDescription)
	// Just make sure the expirationDate is within a minute of the expected
	if claims.ExpiresAt > (expectedClaims.ExpiresAt + 60) {
		t.Errorf("ExpiresAt was greater than expected range for %q, expected %d, got %d",
			testDescription, expectedClaims.ExpiresAt, claims.ExpiresAt)
	}
	if claims.ExpiresAt < (expectedClaims.ExpiresAt - 60) {
		t.Errorf("ExpiresAt was less than expected range for %q, expected %d, got %d",
			testDescription, expectedClaims.ExpiresAt, claims.ExpiresAt)
	}
	return nil
}
