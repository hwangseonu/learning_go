package controllers

import (
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/functions"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/users/requests"
	"github.com/hwangseonu/goBackend/users/responses"
	"gopkg.in/mgo.v2/bson"
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
	} else {
		functions.Response(res, req, 404, []byte(`{"message": "404 page not found"}`))
	}
}

func (c UserController) signUp(res http.ResponseWriter, req *http.Request) {
	var request requests.SignUpRequest
	err := functions.Request(res, req, &request)

	if err != nil {
		return
	}

	err = models.User{
		Id: bson.NewObjectId(),
		Username: request.Username,
		Password: request.Password,
		Nickname: request.Nickname,
		Email:    request.Email,
	}.Save()

	if err != nil {
		if err.Error() == "user already exists" {
			functions.Response(res, req, 409, []byte(`{"message": "user already exists"}`))
			return
		} else {
			functions.Response(res, req, 500, []byte(`{"message": `+err.Error()+`}`))
			return
		}
	}

	functions.Response(res, req, 201, []byte(`{}`))
	return
}

func (c UserController) getUserData(res http.ResponseWriter, req *http.Request) {
	claims := jwt.AuthRequire(res, req, "access")
	if claims == nil {
		return
	}

	user := new(models.User)
	user.FindByUsername(claims.Identity)

	response := responses.GetUserResponse{Username: user.Username, Nickname: user.Nickname, Email: user.Email}
	b, _ := json.MarshalIndent(response, "", "  ")

	functions.Response(res, req, 200, b)
	return
}