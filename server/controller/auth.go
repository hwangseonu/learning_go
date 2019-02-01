package controller

import (
	"net/http"
	"regexp"
)

type AuthController struct {
	http.Handler
}

func (c *AuthController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/auth$").Match(path) && req.Method == "GET" {
		c.signIn(res, req)
	}
}

func (c AuthController) signIn(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")
	_, err := res.Write([]byte(`{"message": "Hello, World"}`))
	if err != nil {
		println(err.Error())
	}
}
