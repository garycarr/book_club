package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/util"
	"github.com/garycarr/book_club/warehouse"
	"github.com/stretchr/testify/assert"
)

const (
	validUserID          = "userID"
	validUserDisplayName = "Gary"
	validUserEmail       = "gcarr@example.com"
	validUserPassword    = "password"
)

func TestLoginPost(t *testing.T) {
	type testData struct {
		description        string
		expectedError      error
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description:        "Valid email and password",
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
	jwt := "Bearer JWT"
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
		mockUtil := util.MockUtil{}
		if td.expectedHTTPStatus == http.StatusOK {
			user := &common.User{
				ID:       validUserID,
				Email:    td.params["email"],
				Password: td.params["password"], // this would normally be bcrypted from the DB
				// As we are mocking everything it doesn't matter
				DisplayName: "JSmith",
			}
			mockWarehouse.On("GetUserWithEmail", td.params["email"]).
				Return(user, nil)
			mockUtil.On("CheckHashedPassword", user.Password, td.params["password"]).Return(nil)
			mockUtil.On("CreateJSONToken", user).Return(jwt, nil)
		} else if td.expectedHTTPStatus == http.StatusUnauthorized {
			mockWarehouse.On("GetUserWithEmail", td.params["email"]).
				Return(nil, common.ErrLoginUserNotFound)
		}
		a.warehouse = &mockWarehouse
		a.util = &mockUtil

		a.Router.ServeHTTP(responseRecorder, req)
		mockWarehouse.AssertExpectations(t)
		mockUtil.AssertExpectations(t)

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

		// The JWT has it's own tests, just make sure it got returned
		tokenString, ok := jsonResp["token"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			continue
		}
		assert.Equal(t, jwt, tokenString, td.description)
	}
}
