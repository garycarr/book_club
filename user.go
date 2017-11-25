package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/garycarr/book_club/common"
)

// loginPost returns a JSON token if the login was successful
func (a *app) userPost(w http.ResponseWriter, r *http.Request) {
	rr := common.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		a.logrus.WithError(err).Error("Unable to decode body")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	if err := rr.ValidateNewUserRequest(); err != nil {
		a.logrus.WithError(err).Error("Missing validation parameters")
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing parameters: %v", err))
		return
	}
	hashedPassword, err := a.util.CreateHashedPassword(rr.Password)
	if err != nil {
		a.logrus.WithError(err).Error("Unable to hash password")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to hash password"))
		return
	}
	rr.Password = string(hashedPassword)
	user, err := a.warehouse.CreateUser(rr)
	if err != nil {
		if err == common.ErrLoginUserAlreadyExists {
			a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Email %s is already registered", rr.Email))
			return
		}
		a.logrus.WithError(err).Error("Unable to create user")
		a.respondWithError(w, http.StatusInternalServerError, "Error creating the user")
		return
	}
	// Create the JSON token as the login is valid
	jsonToken, err := a.util.CreateJSONToken(user)
	if err != nil {
		a.logrus.WithError(err).Error("Unable to create JSON token")
		a.respondWithError(w, http.StatusInternalServerError, "Unable to create JSON token")
		return
	}
	a.respondWithJSON(w, http.StatusCreated, map[string]string{"token": jsonToken})
}

// userOptions returns the allowed options
func (a *app) userOptions(w http.ResponseWriter, r *http.Request) {
	a.optionsHeaders(w)
}
