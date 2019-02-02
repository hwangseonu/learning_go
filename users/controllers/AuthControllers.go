package controllers

import (
	"context"
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/users/requests"
	"github.com/hwangseonu/goBackend/users/responses"
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
	var request requests.SignInRequest
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

	user := new(models.User)
	err = user.FindByUsername(request.Username)
	if err != nil {
		*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 404))
		res.WriteHeader(404)
		res.Write([]byte(`{}`))
		return
	}

	if user.Password != request.Password {
		*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 401))
		res.WriteHeader(401)
		res.Write([]byte(`{}`))
		return
	}
	access, _ := jwt.GenerateToken("access", user.Username)
	refresh, _ := jwt.GenerateToken("refresh", user.Username)
	response := responses.SignInResponse{access, refresh}

	b, _ := json.MarshalIndent(response, "", "  ")

	*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 200))
	res.WriteHeader(200)
	res.Write(b)
	return
}