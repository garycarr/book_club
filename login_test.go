package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestLoginPost(t *testing.T) {
	type testData struct {
		description        string
		expectedHTTPStatus int
		params             map[string]string
	}

	testTable := []testData{
		testData{
			description:        "Valid username and password",
			expectedHTTPStatus: http.StatusOK,
			params: map[string]string{
				"username": validUser,
				"password": validPassword,
			},
		},
		testData{
			description:        "Invalid username",
			expectedHTTPStatus: http.StatusUnauthorized,
			params: map[string]string{
				"username": "invalidUser",
				"password": validPassword,
			},
		},
		testData{
			description:        "Invalid password",
			expectedHTTPStatus: http.StatusUnauthorized,
			params: map[string]string{
				"username": validUser,
				"password": "invalidPassword",
			},
		},
		testData{
			description:        "Password not present",
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"password": validUser,
			},
		},
		testData{
			description:        "Username not present",
			expectedHTTPStatus: http.StatusBadRequest,
			params: map[string]string{
				"username": validUser,
			},
		},
		testData{
			description:        "No params passed in",
			expectedHTTPStatus: http.StatusBadRequest,
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
		resp := executeRequest(req)
		assert.Equal(t, td.expectedHTTPStatus, resp.Code, td.description)
		if td.expectedHTTPStatus == http.StatusOK {
			jsonResp := map[string]string{}
			if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
				t.Errorf("Unable to decode JSON response for test %q: %v", td.description, err)
				continue
			}
			tokenString, ok := jsonResp["token"]
			if !ok {
				t.Errorf("JWT not found for test %q: %v", td.description, jsonResp)
				continue
			}

			claims := jwt.StandardClaims{}
			_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				t.Errorf("Error parsing jwt for test %q: %v", td.description, err)
				continue
			}
			if claims.IssuedAt == 0 || claims.ExpiresAt == 0 || claims.Issuer == "" {
				t.Errorf("Claims were invalid for test %q: %v", td.description, claims)
				continue
			}
		}
	}
}

func TestValidateCredentials(t *testing.T) {
	type testData struct {
		description  string
		expectedResp bool
		lr           LoginRequest
	}

	testTable := []testData{
		testData{
			description:  "Valid username and password",
			expectedResp: true,
			lr: LoginRequest{
				Username: validUser,
				Password: validPassword,
			},
		},
		testData{
			description:  "Invalid username",
			expectedResp: false,
			lr: LoginRequest{
				Username: "invalid_user",
				Password: validPassword,
			},
		},
		testData{
			description:  "Invalid password",
			expectedResp: false,
			lr: LoginRequest{
				Username: validUser,
				Password: "invalid_password",
			},
		},
	}

	for _, td := range testTable {
		if td.lr.validateCredentials() != td.expectedResp {
			t.Errorf("Got unexpected boolean for test %q: %v", td.description, td.expectedResp)
		}
	}
}
