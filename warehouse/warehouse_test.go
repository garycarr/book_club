package warehouse

import (
	"testing"

	"github.com/garycarr/book_club/common"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestWarehouseCreateUserSuccess(t *testing.T) {
	type testData struct {
		description  string
		expectedUser common.User
		rr           common.RegisterRequest
	}

	testTable := []testData{
		testData{
			description: "Should create user",
			expectedUser: common.User{
				ID:          "uniqueRowID",
				DisplayName: "gcarr",
				Email:       "email@example.com",
			},
			rr: common.RegisterRequest{
				Password:    "1234",
				DisplayName: "gcarr",
				Email:       "email@example.com",
			},
		},
	}
	w := Warehouse{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	w.DB = db
	for _, td := range testTable {
		mockInsertUserQuery(mock, td.expectedUser.ID, td.rr.DisplayName, td.rr.Password, td.rr.Email)
		user, createUserErr := w.CreateUser(td.rr)
		if !assert.Nil(t, createUserErr, td.description) {
			// We did not expect an error here, so move onto the next test
		}
		assert.Equal(t, &td.expectedUser, user, td.description)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestWarehouseCreateUserDuplicateEmail(t *testing.T) {
	w := Warehouse{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	w.DB = db
	rr := common.RegisterRequest{
		Password:    "1234",
		DisplayName: "gcarr",
		Email:       "email@example.com",
	}
	mockInsertUserQuery(mock, "uniqueRowID", rr.DisplayName, rr.Password, rr.Email)
	_, err = w.CreateUser(rr)
	if !assert.Nil(t, err) {
		t.Fatal(err)
	}

	// Create the second user, should error
	mock.ExpectQuery("INSERT INTO user_data \\(display_name, password, email\\) VALUES \\(\\$1\\, \\$2, \\$3\\) RETURNING id").
		WithArgs(rr.DisplayName, rr.Password, rr.Email).WillReturnError(common.ErrLoginUserAlreadyExists)

	_, err = w.CreateUser(rr)
	assert.Equal(t, common.ErrLoginUserAlreadyExists, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestWarehouseGetUserWithEmailSuccess(t *testing.T) {
	type testData struct {
		description  string
		expectedUser common.User
		email        string
	}

	testTable := []testData{
		testData{
			description: "Valid email",
			expectedUser: common.User{
				ID:          "1234",
				DisplayName: "gcarr",
				Email:       "email@example.com",
				Password:    "pass123",
			},
			email: "email@example.com",
		},
	}
	w := Warehouse{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	w.DB = db
	for _, td := range testTable {
		mockSelectUserWithEmailQuery(mock, td.expectedUser.ID, td.expectedUser.DisplayName, td.expectedUser.Password, td.expectedUser.Email)
		user, getUserErr := w.GetUserWithEmail(td.email)
		if !assert.Nil(t, getUserErr, td.description) {
			// We did not expect an error here, so move onto the next test
			continue
		}
		// Check the user is what we expected
		assert.Equal(t, &td.expectedUser, user, td.description)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestWarehouseGetUserWithEmailNotFound(t *testing.T) {
	type testData struct {
		description   string
		expectedError error
		email         string
	}

	testTable := []testData{
		testData{
			description:   "No email found",
			expectedError: common.ErrLoginUserNotFound,
			email:         "missing@email.com",
		},
	}
	w := Warehouse{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	w.DB = db
	for _, td := range testTable {
		mock.ExpectQuery("SELECT id, email, password, display_name FROM user_data WHERE email = \\$1").
			WithArgs(td.email).
			WillReturnError(common.ErrLoginUserNotFound)

		_, getUserErr := w.GetUserWithEmail(td.email)
		assert.Equal(t, td.expectedError, getUserErr, td.description)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func mockInsertUserQuery(m sqlmock.Sqlmock, id, displayName, password, email string) {
	m.ExpectQuery("INSERT INTO user_data \\(display_name, password, email\\) VALUES \\(\\$1\\, \\$2, \\$3\\) RETURNING id").
		WithArgs(displayName, password, email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
}

func mockSelectUserWithEmailQuery(m sqlmock.Sqlmock, id, displayName, password, email string) {
	m.ExpectQuery("SELECT id, email, password, display_name FROM user_data WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "display_name"}).AddRow(id, email, password, displayName))
}
