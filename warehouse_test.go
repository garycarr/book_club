package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanUpUserData(t *testing.T, a *app) {
	db, err := a.openDB()
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM user_data")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWarehouseCreateUser(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		expectedUser  user
		rr            registerRequest
	}

	testTable := []testData{
		testData{
			description:   "Should create user",
			expectedError: nil,
			expectedUser: user{
				// id:       "", The ID is a UUID so we cannot know this until it is created
				displayName: "gcarr",
				email:       "email@example.com",
			},
			rr: registerRequest{
				Password:    "1234",
				DisplayName: "gcarr",
				Email:       "email@example.com",
			},
		},
	}
	a := app{}
	if err := a.loadConfiguration("test_config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	for _, td := range testTable {
		user, err := a.createUser(td.rr)
		assert.Equal(t, td.expectedError, err, td.description)
		if td.expectedError != nil {
			cleanUpUserData(t, &a)
			continue
		}
		// Make sure the user was created
		createdUser, err := a.getUserWithEmail(td.rr.Email)
		if err != nil {
			t.Fatal(err)
		}
		td.expectedUser.id = createdUser.id
		assert.Equal(t, &td.expectedUser, user, td.description)
		cleanUpUserData(t, &a)
	}
}

func TestWarehouseCreateUserDuplicate(t *testing.T) {
	a := app{}
	if err := a.loadConfiguration("test_config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	rr := registerRequest{
		Password:    "1234",
		DisplayName: "gcarr",
		Email:       "email@example.com",
	}
	_, err := a.createUser(rr)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	// Create the second user, should error
	_, err = a.createUser(rr)
	assert.Equal(t, errLoginUserAlreadyExists, err)
}

func TestWarehouseGetUserWithEmail(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		expectedUser  user
		email         string
	}

	testTable := []testData{
		testData{
			description:   "Valid email",
			expectedError: nil,
			expectedUser: user{
				id:          "1234",
				displayName: "gcarr",
				email:       "email@example.com",
			},
			email: "email@example.com",
		},
		testData{
			description:   "No email found",
			expectedError: errLoginUserNotFound,
			expectedUser:  user{},
			email:         "missing@email.com",
		},
	}
	a := app{}
	if err := a.loadConfiguration("test_config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	for _, td := range testTable {
		if td.expectedError == nil {
			// We need to put the user in the DB
			createdUser := testCreateUser(t, &a, user{
				email:       td.expectedUser.email,
				password:    "genericPassword",
				displayName: td.expectedUser.displayName,
			}, td.description)
			td.expectedUser.id = createdUser.id
			td.expectedUser.password = createdUser.password
		}
		user, err := a.getUserWithEmail(td.email)
		if !assert.Equal(t, td.expectedError, err, td.description) {
			cleanUpUserData(t, &a)
			continue
		}
		if td.expectedError != nil {
			// We expected an error and it happened and was correct, so move to the next test
			cleanUpUserData(t, &a)
			continue
		}
		// Check the user is what we expected
		assert.Equal(t, &td.expectedUser, user, td.description)
		cleanUpUserData(t, &a)
	}
}
