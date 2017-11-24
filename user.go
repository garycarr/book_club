package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// registerRequest is the information needed to register a new user
type registerRequest struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

// loginPost returns a JSON token if the login was successful
func (a *app) userPost(w http.ResponseWriter, r *http.Request) {
	rr := registerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		a.logrus.WithError(err).Error("Unable to decode body")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	if err := rr.validateNewUserRequest(); err != nil {
		a.logrus.WithError(err).Error("Missing validation parameters")
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing parameters: %v", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcryptCost)
	if err != nil {
		a.logrus.WithError(err).Error("Unable to hash password")
		a.respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to hash password: %v", err))
		return
	}
	rr.Password = string(hashedPassword)
	user, err := a.createUser(rr)
	if err != nil {
		if err == errLoginUserAlreadyExists {
			a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Email %s is already registered", rr.Email))
			return
		}
		a.logrus.WithError(err).Error("Unable to create user")
		a.respondWithError(w, http.StatusInternalServerError, "Error creating the user")
		return
	}

	// Create the JSON token as the login is valid
	jsonToken, err := a.createJSONToken(user)
	if err != nil {
		a.logrus.WithError(err).Error("Unable to create JSON token")
		a.respondWithError(w, http.StatusInternalServerError, "Unable to create JSON token")
		return
	}
	a.respondWithJSON(w, http.StatusCreated, map[string]string{"token": jsonToken})
}

// userOptions returns the allowed options
func (a *app) userOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

func (rr *registerRequest) validateNewUserRequest() error {
	var missingFields string
	if rr.DisplayName == "" {
		missingFields = "displayName,"
	}
	if rr.Password == "" {
		if missingFields != "" {
			missingFields = fmt.Sprintf("%s ", missingFields)
		}
		missingFields += "password,"
	}
	if rr.Email == "" {
		if missingFields != "" {
			missingFields = fmt.Sprintf("%s ", missingFields)
		}
		missingFields += "email,"
	}
	if missingFields != "" {
		missingFields = strings.TrimRight(missingFields, ",")
		return fmt.Errorf(fmt.Sprintf("%s %s", errNewUserMissingFields, missingFields))
	}
	return nil
}
