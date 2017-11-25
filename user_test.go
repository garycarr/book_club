package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/util"
	"github.com/garycarr/book_club/warehouse"
	"github.com/stretchr/testify/assert"
)

func TestUserPost(t *testing.T) {
	type testData struct {
		description        string
		expectedError      error
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description:        "Valid login request",
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
	jwt := "Bearer JWT"
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
			createrUser := &common.User{
				DisplayName: td.params["displayName"],
				Email:       td.params["email"],
			}
			// The bcrypt does not matter here as everything is mocked
			mockUtil.On("CreateHashedPassword", td.params["password"]).Return(td.params["password"], nil)
			mockWarehouse.On("CreateUser", common.RegisterRequest{
				DisplayName: td.params["displayName"],
				Email:       td.params["email"],
				Password:    td.params["password"],
			}).Return(createrUser, nil)
			mockUtil.On("CreateJSONToken", createrUser).Return(jwt, nil)
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

		// The JWT has it's own tests, just make sure it got returned
		tokenString, ok := jsonResp["token"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			continue
		}
		assert.Equal(t, jwt, tokenString, td.description)
	}
}
