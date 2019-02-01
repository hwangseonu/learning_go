package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type AuthController struct {
	http.Handler
}

func (c *AuthController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/auth$").Match(path) && req.Method == "POST" {
		c.signIn(res, req)
	}
}

func (c AuthController) signIn(res http.ResponseWriter, req *http.Request) {
	//TODO
}

type UserController struct {
	http.Handler
}

func (c *UserController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/users$").Match(path) && req.Method == "POST" {
		c.signUp(res, req)
	}
}

func (c UserController) signUp(res http.ResponseWriter, req *http.Request) {
	var request SignUpRequest
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(`{}`))
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte(`{}`))
		return
	}

	res.WriteHeader(201)
	res.Write([]byte(`{}`))
	return
}