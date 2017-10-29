package main

const jwtSecret = "Environmentalise"
const jwtIssuer = "Me"
const port = ":8080"

func main() {
	a := app{}
	a.initialize()
	a.run()
}
