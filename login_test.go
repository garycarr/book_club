package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const (
	validUserDisplayName = "Gary"
	validUserEmail       = "gcarr@example.com"
	validUserPassword    = "password"
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
			description: "Valid email and password",
			expectedClaims: customJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
					Issuer:    jwtIssuer,
				},
				DisplayName: validUserDisplayName,
			},
			expectedHTTPStatus: http.StatusOK,
			params: map[string]string{
				"email":    validUserEmail,
				"password": validUserPassword,
			},
		},
		testData{
			description:        "Wrong email",
			expectedError:      errLoginUserNotFound,
			expectedHTTPStatus: http.StatusUnauthorized,
			params: map[string]string{
				"email":    "invalidUserEmail",
				"password": validUserPassword,
			},
		},
		testData{
			description:        "Wrong password",
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedError:      errLoginUserNotFound,
			params: map[string]string{
				"email":    validUserEmail,
				"password": "invalidUserPassword",
			},
		},
		testData{
			description:        "Email not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginEmailNotPresent,
			params: map[string]string{
				"password": validUserEmail,
			},
		},
		testData{
			description:        "Password not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginPasswordNotPresent,
			params: map[string]string{
				"email": validUserEmail,
			},
		},
		testData{
			description:        "No params passed in",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      errLoginEmailAndPasswordNotPresent,
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
			createdUser := testCreateUser(t, a, user{
				email:       td.params["email"],
				password:    td.params["password"],
				displayName: td.expectedClaims.DisplayName,
			}, td.description)
			td.expectedClaims.Id = createdUser.id
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
