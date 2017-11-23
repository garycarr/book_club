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

func TestUserPost(t *testing.T) {
	type testData struct {
		description        string
		expectedClaims     customJWTClaims
		expectedError      error
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description: "Valid login request",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: "user1",
			},
			expectedHTTPStatus: http.StatusCreated,
			params: map[string]string{
				"email":    "user1@example.com",
				"username": "user1",
				"password": "user1Pass",
			},
		},
		testData{
			description: "Missing username",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: "user1",
			},
			expectedError:      fmt.Errorf("%s username", errNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"email": "user1@example.com",
				// "username": "user1",
				"password": "user1Pass",
			},
		},
		testData{
			description: "Missing email",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: "user1",
			},
			expectedError:      fmt.Errorf("%s email", errNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				// "email": "user1@example.com",
				"username": "user1",
				"password": "user1Pass",
			},
		},
		testData{
			description: "Missing password",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: "user1",
			},
			expectedError:      fmt.Errorf("%s password", errNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"email":    "user1@example.com",
				"username": "user1",
				// "password": "user1Pass",
			},
		},
		testData{
			description: "Missing everything",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				Username: "user1",
			},
			expectedError:      fmt.Errorf("%s username, password, email", errNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params:             map[string]string{
			// "email": "user1@example.com",
			// "username": "user1",
			// "password": "user1Pass",
			},
		},
	}
	for _, td := range testTable {
		params, err := json.Marshal(td.params)
		if err != nil {
			t.Fatalf("Error marshalling for test %q: %v", td.description, err)
		}
		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(params))
		if err != nil {
			t.Fatalf("Error creating new request for test %q: %v", td.description, err)
		}
		a, responseRecorder := setupTest(req)
		defer cleanUpUserData(t, a)
		a.Router.ServeHTTP(responseRecorder, req)
		if !assert.Equal(t, td.expectedHTTPStatus, responseRecorder.Code, td.description) {
			// We got a different status code than expected
			continue
		}

		jsonResp := map[string]string{}
		if err = json.NewDecoder(responseRecorder.Body).Decode(&jsonResp); err != nil {
			t.Errorf("Unable to decode JSON response for test %q: %v", td.description, err)
			continue
		}
		// Check error message
		if td.expectedHTTPStatus != http.StatusCreated {
			assert.Contains(t, jsonResp["error"], td.expectedError.Error(), td.description)
			continue
		}

		// Get the created user so we can check the ID
		createdUser, err := a.getUserWithUsername(td.params["username"])
		if err != nil {
			t.Errorf("Unable to get user for test %q: %v", td.description, err)
		}
		td.expectedClaims.Id = createdUser.id

		// JWT tests
		tokenString, ok := jsonResp["token"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			continue
		}

		if err = checkJWT(t, td.expectedClaims, tokenString, td.description); err != nil {
			t.Error(err)
			continue
		}
	}
}
