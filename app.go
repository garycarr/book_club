package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type app struct {
	Router *mux.Router
	conf   *config
}

type config struct {
	Database struct {
		DBName   string `json:"db_name"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"database"`
}

func (a *app) run() {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *app) openDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		a.conf.Database.Username, a.conf.Database.Password, a.conf.Database.Host, a.conf.Database.DBName)
	return sql.Open("postgres", connectionString)
}

func (a *app) initialize() {
	err := a.loadConfiguration("config.json")
	if err != nil {
		panic(err)
	}

	db, err := a.openDB()
	defer db.Close()
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *app) initializeRoutes() {
	a.Router.HandleFunc("/login", a.loginPost).Methods("POST")
	a.Router.HandleFunc("/login", a.loginOptions).Methods("OPTIONS")
}

func (a *app) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *app) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *app) loadConfiguration(file string) error {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&a.conf); err != nil {
		return err
	}
	return nil
}
