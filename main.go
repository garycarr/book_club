package main

import "time"

const bcryptCost = 10
const jwtSecret = "Environmentalise"
const jwtIssuer = "Me"
const port = ":8080"
const jwtExpiration = time.Duration(1 * time.Hour)

func main() {
	a := app{}
	a.initialize("config.json")
	a.logrus.Info("Starting the app")
	a.run()
}
