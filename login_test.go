package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/warehouse"
	"github.com/stretchr/testify/assert"
)

const (
	validUserID          = "userID"
	validUserDisplayName = "Gary"
	validUserEmail       = "gcarr@example.com"
	validUserPassword    = "password"
	bcryptPassword       = "$2a$10$O02CA3UNrCC12JU66PIWvuI4oJceIUSB7Y/FX1x9ujzXsutCoeIMS"
)

func TestLoginPost(t *testing.T) {
	type testData struct {
		description        string
		expectedClaims     common.CustomJWTClaims
		expectedError      error
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description: "Valid email and password",
			expectedClaims: common.CustomJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(common.JWTExpiration).Unix(),
					Issuer:    common.JWTIssuer,
					Id:        validUserID,
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
			expectedError:      common.ErrLoginUserNotFound,
			expectedHTTPStatus: http.StatusUnauthorized,
			params: map[string]string{
				"email":    "invalidUserEmail",
				"password": validUserPassword,
			},
		},
		testData{
			description:        "Wrong password",
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedError:      common.ErrLoginUserNotFound,
			params: map[string]string{
				"email":    validUserEmail,
				"password": "invalidUserPassword",
			},
		},
		testData{
			description:        "Email not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      common.ErrLoginEmailNotPresent,
			params: map[string]string{
				"password": validUserEmail,
			},
		},
		testData{
			description:        "Password not present",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      common.ErrLoginPasswordNotPresent,
			params: map[string]string{
				"email": validUserEmail,
			},
		},
		testData{
			description:        "No params passed in",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedError:      common.ErrLoginEmailAndPasswordNotPresent,
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
		mockWarehouse := warehouse.MockWarehouse{}
		if td.expectedHTTPStatus == http.StatusOK {
			mockWarehouse.On("GetUserWithEmail", td.params["email"]).
				Return(&common.User{
					ID:          validUserID,
					Email:       td.params["email"],
					Password:    bcryptPassword,
					DisplayName: td.expectedClaims.DisplayName,
				}, nil)
		} else if td.expectedHTTPStatus == http.StatusUnauthorized {
			mockWarehouse.On("GetUserWithEmail", td.params["email"]).
				Return(nil, common.ErrLoginUserNotFound)
		}
		a.warehouse = &mockWarehouse

		a.Router.ServeHTTP(responseRecorder, req)
		mockWarehouse.AssertExpectations(t)
		if !assert.Equal(t, td.expectedHTTPStatus, responseRecorder.Code, td.description) {
			// We got a different status code than expected
			continue
		}

		jsonResp := map[string]string{}
		if err = json.NewDecoder(responseRecorder.Body).Decode(&jsonResp); err != nil {
			t.Errorf("Unable to decode JSON response for test %q: %v", td.description, err)
			continue
		}
		if td.expectedHTTPStatus != http.StatusOK {
			assert.Contains(t, jsonResp["error"], td.expectedError.Error(), td.description)
			// We we were expecting an error there is nothing else to check
			continue
		}

		// JWT tests
		tokenString, ok := jsonResp["token"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			continue
		}
		if err = checkJWT(t, td.expectedClaims, tokenString, td.description); err != nil {
			t.Error(err)
		}
	}
}
