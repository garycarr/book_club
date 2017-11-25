package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/garycarr/book_club/util"
	"github.com/garycarr/book_club/warehouse"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type app struct {
	conf             *config
	connectionString string
	logrus           *logrus.Logger
	util             util.UtilIn
	Router           *mux.Router
	warehouse        warehouse.WarehouseIn
}

type config struct {
	Port     string `json:"port"`
	Database struct {
		DBName   string `json:"db_name"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"database"`
}

func (a *app) run() {
	a.logrus.Fatal(http.ListenAndServe(a.conf.Port, a.Router))
}

func (a *app) initialize(configFile string) {
	a.logrus = logrus.New()
	a.logrus.Formatter = &logrus.JSONFormatter{}
	err := a.loadConfiguration(configFile)
	if err != nil {
		a.logrus.WithError(err).Fatal("Error loading config")
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		a.conf.Database.Username, a.conf.Database.Password, a.conf.Database.Host, a.conf.Database.DBName)

	wh, err := warehouse.NewWarehouse(connectionString, a.logrus)
	if err != nil {
		a.logrus.WithError(err).Fatal("Error creating warehouse")
	}
	a.warehouse = wh
	a.util = util.NewUtil()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *app) initializeRoutes() {
	a.Router.HandleFunc("/login", a.loginPost).Methods(http.MethodPost)
	a.Router.HandleFunc("/login", a.loginOptions).Methods(http.MethodOptions)

	a.Router.HandleFunc("/user", a.userPost).Methods(http.MethodPost)
	a.Router.HandleFunc("/user", a.userOptions).Methods(http.MethodOptions)

	authMiddleware := alice.New(a.authMiddleware)
	a.Router.Handle("/homepage", authMiddleware.ThenFunc(a.homePageGet)).Methods(http.MethodGet)
	a.Router.Handle("/homepage", authMiddleware.ThenFunc(a.homePageOptions)).Methods(http.MethodOptions)
}

func (a *app) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *app) authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := a.util.CheckJSONToken(r.Header.Get("Authorization")); err != nil {
			a.logrus.WithError(err).Debug("Ivalid JSON token. Redirecting user to homepage")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
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

func (a *app) optionsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}
