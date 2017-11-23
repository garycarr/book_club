package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	validUser     = "gcarr"
	validPassword = "password"
)

func TestLoginPost(t *testing.T) {
	type testData struct {
		description        string
		expectedClaims     customJWTClaims
		expectedError      error
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description: "Valid username and password",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: validUser,
			},
			expectedHTTPStatus: http.StatusOK,
			params: map[string]string{
				"username": validUser,
				"password": validPassword,
			},
		},
		testData{
			description:        "Wrong username",
			expectedError:      errLoginUserNotFound,
			expectedHTTPStatus: http.StatusUnauthorized,
			params: map[string]string{
				"username": "invalidUser",
				"password": validPassword,
			},
		},
		testData{
			description:        "Wrong password",
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedError:      errLoginUserNotFound,
			params: map[string]string{
				"username": validUser,
				"password": "invalidPassword",
			},
		},
		testData{
			description:        "Username not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginUsernameNotPresent,
			params: map[string]string{
				"password": validUser,
			},
		},
		testData{
			description:        "Password not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginPasswordNotPresent,
			params: map[string]string{
				"username": validUser,
			},
		},
		testData{
			description:        "No params passed in",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginUsernameAndPasswordNotPresent,
		},
	}

	for _, td := range testTable {
		params, err := json.Marshal(td.params)
		if err != nil {
			t.Fatalf("Error marshalling for test %q: %v", td.description, err)
		}
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(params))
		if err != nil {
			t.Fatalf("Error creating new request for test %q: %v", td.description, err)
		}
		a, responseRecorder := setupTest(req)
		if td.expectedHTTPStatus == http.StatusOK {
			// We need to put the user in the DB
			td.expectedClaims.Id = testCreateUser(t, a, user{
				email:    "generic@email.com",
				password: td.params["password"],
				username: td.params["username"],
			}, td.description)
		}
		a.Router.ServeHTTP(responseRecorder, req)
		if !assert.Equal(t, td.expectedHTTPStatus, responseRecorder.Code, td.description) {
			// We got a different status code than expected
			cleanUpUserData(t, a)
			continue
		}

		jsonResp := map[string]string{}
		if err = json.NewDecoder(responseRecorder.Body).Decode(&jsonResp); err != nil {
			t.Errorf("Unable to decode JSON response for test %q: %v", td.description, err)
			cleanUpUserData(t, a)
			continue
		}
		// Check error message
		if td.expectedHTTPStatus != http.StatusOK {
			assert.Contains(t, jsonResp["error"], td.expectedError.Error(), td.description)
			cleanUpUserData(t, a)
			continue
		}

		// JWT tests
		tokenString, ok := jsonResp["token"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			cleanUpUserData(t, a)
			continue
		}

		if err = checkJWT(t, td.expectedClaims, tokenString, td.description); err != nil {
			t.Error(err)
			cleanUpUserData(t, a)
			continue
		}
		cleanUpUserData(t, a)
	}
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
