package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type app struct {
	Router *mux.Router
}

func (a *app) run() {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *app) initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *app) initializeRoutes() {
	a.Router.HandleFunc("/login", loginPost).Methods("POST")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
