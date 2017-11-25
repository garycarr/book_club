package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/util"
	"github.com/stretchr/testify/assert"
)

type authTestData struct {
	description        string
	expectedHTTPStatus int
	jwToken            string
	validJWT           bool
}

func authJWTTestTable(t *testing.T) []authTestData {
	util := util.NewUtil()
	jwToken, err := util.CreateJSONToken(&common.User{
		ID:          "123",
		Email:       "email@example.com",
		DisplayName: "JSmith",
	})
	if err != nil {
		t.Fatal(err)
	}

	return []authTestData{
		authTestData{
			description:        "Valid JWT",
			expectedHTTPStatus: http.StatusOK,
			jwToken:            fmt.Sprintf("Bearer %s", jwToken),
			validJWT:           true,
		},
		authTestData{
			description:        "Invalid JWT",
			expectedHTTPStatus: http.StatusUnauthorized,
			jwToken:            fmt.Sprintf("%s", jwToken),
			validJWT:           false,
		},
		authTestData{
			description:        "No JWT",
			expectedHTTPStatus: http.StatusUnauthorized,
			jwToken:            "",
			validJWT:           false,
		},
	}
}

func setupTest(req *http.Request) (*app, *httptest.ResponseRecorder) {
	rr := httptest.NewRecorder()
	a := app{}
	a.initialize("test_config.json")
	return &a, rr
}

func authEndpointTests(t *testing.T, path string) {
	for _, td := range authJWTTestTable(t) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			t.Fatalf("Error creating new request for test %q: %v", td.description, err)
		}
		req.Header.Add("Authorization", td.jwToken)
		a, responseRecorder := setupTest(req)
		a.Router.ServeHTTP(responseRecorder, req)
		assert.Equal(t, td.expectedHTTPStatus, responseRecorder.Code, td.description)
	}
}
