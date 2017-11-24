package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// loginRequest is the data needed to make a login
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginPost returns a JSON token if the login was successful
func (a *app) loginPost(w http.ResponseWriter, r *http.Request) {
	login := loginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		a.logrus.WithError(err).Error("Unable to decode body")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	if err := login.validateRequest(); err != nil {
		a.logrus.WithError(err).Error("Missing validation parameters")
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing parameters: %v", err))
		return
	}
	user, err := a.validateCredentials(login)
	if err != nil {
		if err == errLoginUserNotFound {
			a.logrus.Debug("Incorrect password given")
			a.respondWithError(w, http.StatusUnauthorized, err.Error())
		} else {
			a.logrus.WithError(err).Error("Incorrect login details")
			a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking user credentials"))
		}
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

func (lr *loginRequest) validateRequest() error {
	if lr.Email == "" && lr.Password == "" {
		return errLoginEmailAndPasswordNotPresent
	}
	if lr.Email == "" {
		return errLoginEmailNotPresent
	}
	if lr.Password == "" {
		return errLoginPasswordNotPresent
	}
	return nil
}

func (a *app) validateCredentials(lr loginRequest) (*user, error) {
	user, err := a.getUserWithEmail(lr.Email)
	if err != nil {
		return nil, err
	}
	// Make sure the password is valid
	if err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(lr.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return nil, fmt.Errorf("An error occurred checking the password")
		}
		return nil, errLoginUserNotFound
	}
	return user, nil
}
