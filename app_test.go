package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
	assert.Equal(t, expectedClaims.DisplayName, claims.DisplayName, testDescription)
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

func testCreateUser(t *testing.T, a *app, u user, testDescription string) *user {
	// We need to put the user in the DB
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcryptCost)
	if err != nil {
		t.Fatalf("Could not create password for %q, err: %q", testDescription, err.Error())
	}
	createdUser, err := a.createUser(registerRequest{
		Email:       u.email,
		Password:    string(hashedPassword),
		DisplayName: u.displayName,
	})
	if err != nil {
		t.Fatalf("Error inserting into DB for %q: %q", testDescription, err.Error())
	}

	// For testing convenience return the hashedPassword
	createdUser.password = string(hashedPassword)
	return createdUser
}
