package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/util"
	"github.com/garycarr/book_club/warehouse"
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
				DisplayName: "user1",
			},
			expectedHTTPStatus: http.StatusCreated,
			params: map[string]string{
				"email":       "user1@example.com",
				"displayName": "user1",
				"password":    "user1Pass",
			},
		},
		testData{
			description:        "Missing displayName",
			expectedError:      fmt.Errorf("%s displayName", common.ErrNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"email": "user2@example.com",
				// "displayName": "user2",
				"password": "user2Pass",
			},
		},
		testData{
			description:        "Missing email",
			expectedError:      fmt.Errorf("%s email", common.ErrNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				// "email": "user3@example.com",
				"displayName": "user3",
				"password":    "user3Pass",
			},
		},
		testData{
			description:        "Missing password",
			expectedError:      fmt.Errorf("%s password", common.ErrNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"email":       "user4@example.com",
				"displayName": "user4",
				// "password": "user4Pass",
			},
		},
		testData{
			description:        "Missing everything",
			expectedError:      fmt.Errorf("%s displayName, password, email", common.ErrNewUserMissingFields),
			expectedHTTPStatus: http.StatusBadRequest,
			params:             map[string]string{
			// "email": "user5@example.com",
			// "displayName": "user5",
			// "password": "user5Pass",
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
		mockUtil := util.MockUtil{}
		mockWarehouse := warehouse.MockWarehouse{}
		if td.expectedError == nil {
			mockUtil.On("GetCryptedPassword", td.params["password"]).Return(bcryptPassword, nil)
			mockWarehouse.On("CreateUser", common.RegisterRequest{
				DisplayName: td.params["displayName"],
				Email:       td.params["email"],
				Password:    bcryptPassword,
			}).Return(&common.User{
				DisplayName: td.params["displayName"],
				Email:       td.params["email"],
			}, nil)
		}
		a.util = &mockUtil
		a.warehouse = &mockWarehouse

		a.Router.ServeHTTP(responseRecorder, req)
		mockUtil.AssertExpectations(t)
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
		if td.expectedHTTPStatus != http.StatusCreated {
			assert.Contains(t, jsonResp["error"], td.expectedError.Error(), td.description)
			// We were expecting an error, so move onto the next test
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
			continue
		}
	}
}
