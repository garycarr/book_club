package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	a := app{}
	a.initialize("config.json")
	defer a.warehouse.Close()
	a.logrus.WithFields(logrus.Fields{
		"dbHost":   a.conf.Database.Host,
		"dbName":   a.conf.Database.DBName,
		"logLevel": a.logrus.Level,
	}).Info("Starting the app")
	a.run()
}
