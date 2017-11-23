package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// loginRequest is the information needed to make a login
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type customJWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// loginPost returns a JSON token if the login was successful
func (a *app) loginPost(w http.ResponseWriter, r *http.Request) {
	login := loginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	if err := login.validateRequest(); err != nil {
		a.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Missing parameters: %v", err))
		return
	}
	user, err := a.validateCredentials(login)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == errLoginUserNotFound {
			statusCode = http.StatusUnauthorized
		}
		a.respondWithError(w, statusCode, fmt.Sprintf("Error getting user: %s", err.Error()))
		return
	}
	// Create the JSON token as the login is valid
	claims := &customJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
			Issuer:    jwtIssuer,
			IssuedAt:  time.Now().Unix(),
			Id:        user.id,
		},
		Username: user.username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		a.respondWithError(w, http.StatusUnauthorized, "Unable to create JSON token")
		return
	}
	a.respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// loginOptions returns the allowed options
func (a *app) loginOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

func (lr *loginRequest) validateRequest() error {
	if lr.Username == "" && lr.Password == "" {
		return errLoginUsernameAndPasswordNotPresent
	}
	if lr.Username == "" {
		return errLoginUsernameNotPresent
	}
	if lr.Password == "" {
		return errLoginPasswordNotPresent
	}
	return nil
}

func (a *app) validateCredentials(lr loginRequest) (*user, error) {
	// Need to hash password
	user, err := a.getUserWithPassword(lr)
	if err != nil {
		return nil, err
	}
	return user, nil
}
