package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

const bcryptCost = 10
const jwtSecret = "Environmentalise"
const jwtIssuer = "Me"
const port = ":8080"
const jwtExpiration = time.Duration(1 * time.Hour)

func main() {
	a := app{}
	a.initialize("config.json")
	defer a.warehouse.DB.Close()
	a.logrus.WithFields(logrus.Fields{
		"dbHost":   a.conf.Database.Host,
		"dbName":   a.conf.Database.DBName,
		"logLevel": a.logrus.Level,
	}).Info("Starting the app")
	a.run()
}
