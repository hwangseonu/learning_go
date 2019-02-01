package users

import (
	"context"
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
		*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 400))
		res.WriteHeader(400)
		res.Write([]byte(`{}`))
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 400))
		res.WriteHeader(400)
		res.Write([]byte(`{}`))
		return
	}

	err = User{
		request.Username,
		request.Password,
		request.Nickname,
		request.Email,
	}.Save()

	if err != nil {
		if err.Error() == "user already exists" {
			*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 409))
			res.WriteHeader(409)
			res.Write([]byte(`{}`))
			return
		} else {
			*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 500))
			res.WriteHeader(500)
			res.Write([]byte(`{"message": `+err.Error()+`}`))
			return
		}
	}

	*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 201))
	res.WriteHeader(201)
	res.Write([]byte(`{}`))
	return
}