package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/garycarr/book_club/common"
	"golang.org/x/crypto/bcrypt"
)

// loginPost returns a JSON token if the login was successful
func (a *app) loginPost(w http.ResponseWriter, r *http.Request) {
	login := common.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		a.logrus.WithError(err).Error("Unable to decode body")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	if err := login.ValidateRequest(); err != nil {
		a.logrus.WithError(err).Error("Missing validation parameters")
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing parameters: %v", err))
		return
	}
	user, err := a.validateCredentials(login)
	if err != nil {
		if err == common.ErrLoginUserNotFound {
			a.logrus.Debug("Incorrect password given")
			a.respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		a.logrus.WithError(err).Error("Incorrffffect login details")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking user credentials"))
		return
	}
	// Create the JSON token as the login is valid
	jsonToken, err := a.createJSONToken(user)
	if err != nil {
		a.logrus.WithError(err).Error("Unable to create JSON token")
		a.respondWithError(w, http.StatusInternalServerError, "Unable to create JSON token")
		return
	}
	a.respondWithJSON(w, http.StatusOK, map[string]string{"token": jsonToken})
}

// loginOptions returns the allowed options
func (a *app) loginOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

func (a *app) validateCredentials(lr common.LoginRequest) (*common.User, error) {
	user, err := a.warehouse.GetUserWithEmail(lr.Email)
	if err != nil {
		return nil, err
	}
	// Make sure the password is valid
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("An error occurred checking the password: %s", err.Error())
		}
		return nil, common.ErrLoginUserNotFound
	}
	return user, nil
}
