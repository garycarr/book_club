package main

import (
	"net/http"
)

// homePageGet returns the homepage data
func (a *app) homePageGet(w http.ResponseWriter, r *http.Request) {
	a.respondWithJSON(w, http.StatusOK, map[string]string{"salrightman": "boom"})
}

// homePageOptions returns the allowed options
func (a *app) homePageOptions(w http.ResponseWriter, r *http.Request) {
	a.optionsHeaders(w)
}
