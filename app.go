package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/garycarr/book_club/common"
	"github.com/garycarr/book_club/util"
	"github.com/garycarr/book_club/warehouse"
	"github.com/gorilla/mux"
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
	Database struct {
		DBName   string `json:"db_name"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"database"`
}

type customJWTClaims struct {
	DisplayName string `json:"displayName"`
	jwt.StandardClaims
}

func (a *app) run() {
	a.logrus.Fatal(http.ListenAndServe(port, a.Router))
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
	a.Router.HandleFunc("/login", a.loginPost).Methods("POST")
	a.Router.HandleFunc("/login", a.loginOptions).Methods("OPTIONS")

	a.Router.HandleFunc("/user", a.userPost).Methods("POST")
	a.Router.HandleFunc("/user", a.userOptions).Methods("OPTIONS")
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

func (a *app) createJSONToken(u *common.User) (string, error) {
	// Create the JSON token as the login is valid
	claims := &customJWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
			Issuer:    jwtIssuer,
			IssuedAt:  time.Now().Unix(),
			Id:        u.ID,
		},
		DisplayName: u.DisplayName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
