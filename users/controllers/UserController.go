package controllers

import (
	"context"
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/users/requests"
	"io/ioutil"
	"net/http"
	"regexp"
)

type UserController struct {
	http.Handler
}

func (c *UserController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/users$").Match(path) && req.Method == "POST" {
		c.signUp(res, req)
	} else if regexp.MustCompile("^/users$").Match(path) && req.Method == "GET" {
		c.getUserData(res, req)
	}
}

func (c UserController) signUp(res http.ResponseWriter, req *http.Request) {
	var request requests.SignUpRequest
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

	err = models.User{
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

func (c UserController) getUserData(res http.ResponseWriter, req *http.Request) {
	claims := jwt.AuthRequire(res, req, "access")
	if claims == nil {
		return
	}
	res.WriteHeader(200)
	res.Write([]byte(`{}`))
	return
}