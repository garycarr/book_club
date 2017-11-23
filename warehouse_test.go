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
				username: "gcarr",
				email:    "email@example.com",
			},
			rr: registerRequest{
				Password: "1234",
				Username: "gcarr",
				Email:    "email@example.com",
			},
		},
	}
	a := app{}
	if err := a.loadConfiguration("config.json"); err != nil {
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
		createdUser, err := a.getUserWithUsername(td.rr.Username)
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
	if err := a.loadConfiguration("config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	rr := registerRequest{
		Password: "1234",
		Username: "gcarr",
		Email:    "email@example.com",
	}
	_, err := a.createUser(rr)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	// Create the second user, should error
	_, err = a.createUser(rr)
	assert.Equal(t, errLoginUserAlreadyExists, err)
}

func TestWarehouseGetUserWithPassword(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		expectedUser  user
		lr            loginRequest
	}

	testTable := []testData{
		testData{
			description:   "Valid username and password",
			expectedError: nil,
			expectedUser: user{
				id:       "1234",
				username: "gcarr",
				email:    "email@example.com",
			},
			lr: loginRequest{
				Password: "pass",
				Username: "gcarr",
			},
		},
		testData{
			description:   "Invalid password",
			expectedError: errLoginUserNotFound,
			expectedUser:  user{},
			lr: loginRequest{
				Password: "invalidPass",
				Username: "gcarr",
			},
		},
	}
	a := app{}
	if err := a.loadConfiguration("config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	for _, td := range testTable {
		if td.expectedError == nil {
			// We need to put the user in the DB
			td.expectedUser.id = testCreateUser(t, &a, user{
				email:    td.expectedUser.email,
				password: td.lr.Password,
				username: td.lr.Username,
			}, td.description)
		}
		user, err := a.getUserWithPassword(td.lr)
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

func TestWarehouseGetUserWithUsername(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		expectedUser  user
		username      string
	}

	testTable := []testData{
		testData{
			description:   "Valid username",
			expectedError: nil,
			expectedUser: user{
				id:       "1234",
				username: "gcarr",
				email:    "email@example.com",
			},
			username: "gcarr",
		},
		testData{
			description:   "Invalid password",
			expectedError: errLoginUserNotFound,
			expectedUser:  user{},
			username:      "invalidUser",
		},
	}
	a := app{}
	if err := a.loadConfiguration("config.json"); err != nil {
		t.Fatal(err)
	}
	defer cleanUpUserData(t, &a)
	for _, td := range testTable {
		if td.expectedError == nil {
			// We need to put the user in the DB
			td.expectedUser.id = testCreateUser(t, &a, user{
				email:    td.expectedUser.email,
				password: "genericPassword",
				username: td.username,
			}, td.description)
		}
		user, err := a.getUserWithUsername(td.username)
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

func testCreateUser(t *testing.T, a *app, u user, testDescription string) string {
	// We need to put the user in the DB
	createdUser, err := a.createUser(registerRequest{
		Email:    u.email,
		Password: u.password,
		Username: u.username,
	})
	if err != nil {
		t.Fatalf("Error inserting into DB for %q: %q", testDescription, err.Error())
	}
	return createdUser.id
}
