package main

import (
	"database/sql"

	"github.com/lib/pq"
)

type user struct {
	email       string
	id          string
	displayName string
	password    string
}

func (a *app) createUser(rr registerRequest) (*user, error) {
	db, err := a.openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var id string
	sqlStatement := `INSERT INTO user_data (display_name, password, email)
		VALUES ($1, $2, $3)
		RETURNING id`
	if err = db.QueryRow(sqlStatement, rr.DisplayName, rr.Password, rr.Email).Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, errLoginUserAlreadyExists
		}
		return nil, err
	}
	return &user{
		id:          id,
		email:       rr.Email,
		displayName: rr.DisplayName,
	}, nil
}

func (a *app) getUserWithEmail(email string) (*user, error) {
	u := user{}
	db, err := a.openDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	sqlStatement := `SELECT id, email, password, display_name
		FROM user_data
		WHERE email = $1`
	err = db.QueryRow(sqlStatement, email).Scan(&u.id, &u.email, &u.password, &u.displayName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errLoginUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
