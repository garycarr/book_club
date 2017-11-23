package main

import "time"

const jwtSecret = "Environmentalise"
const jwtIssuer = "Me"
const port = ":8080"
const jwtExpiration = time.Duration(1 * time.Hour)

func main() {
	a := app{}
	a.initialize()
	a.run()
}
