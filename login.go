package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Get consts from DB
const (
	validUser     = "gcarr"
	validPassword = "password"
)

// LoginRequest is the information needed to make a login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginPost returns a JSON token if the login was successful
func loginPost(w http.ResponseWriter, r *http.Request) {
	login := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unable to decode request: %v", err))
		return
	}
	var errorMessage string
	if login.Username == "" {
		errorMessage = "Username not present."
	}
	if login.Password == "" {
		errorMessage = fmt.Sprintf("%s Password not Present.", errorMessage)
	}
	if errorMessage != "" {
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}
	if !login.validateCredentials() {
		respondWithError(w, http.StatusUnauthorized, "Username or password not found")
		return
	}

	// Create the JSON token as the login is valid
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		Issuer:    jwtIssuer,
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to create JSON token")
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (lr *LoginRequest) validateCredentials() bool {
	if lr.Username != validUser {
		return false
	}
	if lr.Password != validPassword {
		return false
	}
	return true
}
