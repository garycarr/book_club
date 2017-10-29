package main

import (
	"net/http"
	"net/http/httptest"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a := app{}
	a.initialize()
	a.Router.ServeHTTP(rr, req)
	return rr
}
