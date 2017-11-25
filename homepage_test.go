package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/garycarr/book_club/util"
	"github.com/stretchr/testify/assert"
)

func TestHomePageGetJWTTests(t *testing.T) {
	authEndpointTests(t, "/homepage")
}

func TestHomePageGet(t *testing.T) {
	type testData struct {
		description        string
		expectedError      error
		expectedHTTPStatus int
		expectedResponse   string
	}

	testTable := []testData{
		testData{
			description:        "Successful GET request",
			expectedHTTPStatus: http.StatusOK,
		},
	}
	validJWT := "JWT"
	for _, td := range testTable {
		req, err := http.NewRequest(http.MethodGet, "/homepage", nil)
		if err != nil {
			t.Fatalf("Error creating new request for test %q: %v", td.description, err)
		}
		req.Header.Add("Authorization", validJWT)
		a, responseRecorder := setupTest(req)
		mockUtil := util.MockUtil{}
		mockUtil.On("CheckJSONToken", validJWT).Return(nil)
		a.util = &mockUtil
		a.Router.ServeHTTP(responseRecorder, req)
		assert.Equal(t, td.expectedHTTPStatus, responseRecorder.Code, td.description)

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

		// Stubbed :)
		salrightman, ok := jsonResp["salrightman"]
		if !ok {
			t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
			continue
		}
		assert.Equal(t, salrightman, "boom", td.description)
	}
}
