package util

import (
	"fmt"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garycarr/book_club/common"
	"github.com/stretchr/testify/assert"
)

func TestCreateJSONToken(t *testing.T) {
	type testData struct {
		description    string
		expectedClaims customJWTClaims
		expectedError  error
		user           *common.User
	}

	testTable := []testData{
		testData{
			description: "should create token",
			user: &common.User{
				DisplayName: "Bob",
				ID:          "abc123",
			},
		},
	}
	u := NewUtil()
	for _, td := range testTable {
		jwToken, err := u.CreateJSONToken(td.user)
		if err != nil {
			t.Errorf("Unexpected err for %q: %v", td.description, td.expectedError)
		}
		claims := customJWTClaims{}
		_, err = jwt.ParseWithClaims(jwToken, &claims, func(jwToken *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil {
			t.Errorf("Error parsing jwt for %q: %v", td.description, err)
			continue
		}
		//
		assert.Equal(t, td.user.ID, claims.Id, td.description)
		assert.Equal(t, jwtIssuer, claims.Issuer, td.description)
		assert.Equal(t, td.user.DisplayName, claims.DisplayName, td.description)
		// Just make sure the expirationDate is within a minute of the expected
		if claims.ExpiresAt > (time.Now().Add(jwtExpiration).Add(time.Hour).Unix()) {
			t.Errorf("ExpiresAt was greater than expected range for %q, got %d",
				td.description, claims.ExpiresAt)
		}
		if claims.ExpiresAt < (time.Now().Add(jwtExpiration).Add(-time.Hour).Unix()) {
			t.Errorf("ExpiresAt was less than expected range for %q, got %d",
				td.description, claims.ExpiresAt)
		}
	}
}

func TestCheckJSONToken(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		invalidJWT    string
		user          *common.User
	}

	jwtForFailingTests := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkaXNwbGF5TmFtZSI6IkJvYiIsImV4cCI6MTUxMTYyMDk3MywianRpIjoiYWJjMTIzIiwiaWF0IjoxNTExNjE3MzczLCJpc3MiOiJNZSJ9.EM4uTvEpT4bvjsiPwZJGl_qcPRSKh5snzQ_rzMP82gs"
	testTable := []testData{
		testData{
			description: "valid token",
			user: &common.User{
				DisplayName: "Bob",
				ID:          "abc123",
			},
		},
		testData{
			description:   "Missing Bearer",
			expectedError: common.ErrJSONTokenNoBearer,
			invalidJWT:    jwtForFailingTests,
		},
		testData{
			description:   "Invalid With bearer",
			expectedError: jwt.ValidationError{Inner: jwt.ErrSignatureInvalid},
			invalidJWT:    fmt.Sprintf("Bearer %s1", jwtForFailingTests), // The one on the end makes it invalid
		},
	}
	u := NewUtil()
	for _, td := range testTable {
		// Make a valid token - yeah yeah, shouldn't use prod functions for tests
		if td.expectedError == nil {
			jwToken, err := u.CreateJSONToken(td.user)
			jwToken = fmt.Sprintf("Bearer %s", jwToken)
			if err != nil {
				t.Errorf("Unexpected err for %q: %v", td.description, td.expectedError)
			}
			assert.Nil(t, u.CheckJSONToken(jwToken), td.description)
		} else {
			err := u.CheckJSONToken(td.invalidJWT)
			assert.Equal(t, err.Error(), td.expectedError.Error(), td.description)
		}
	}
}
