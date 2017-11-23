package main

import (
	"net/http"
	"net/http/httptest"
)

func setupTest(req *http.Request) (*app, *httptest.ResponseRecorder) {
	rr := httptest.NewRecorder()
	a := app{}
	a.loadConfiguration("config.json")
	a.initialize()
	return &a, rr
}
