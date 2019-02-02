package controllers

import (
	"context"
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/functions"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/users/requests"
	"github.com/hwangseonu/goBackend/users/responses"
	"net/http"
	"regexp"
	"time"
)

type AuthController struct {
	http.Handler
}

func (c *AuthController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/auth$").Match(path) && req.Method == "POST" {
		c.signIn(res, req)
	} else if regexp.MustCompile("^/auth/refresh$").Match(path) && req.Method == "GET" {
		c.refresh(res, req)
	}
	functions.Response(res, req, 404, []byte(`{"message": "404 page not found"}`))
}

func (c AuthController) signIn(res http.ResponseWriter, req *http.Request) {
	var request requests.SignInRequest
	err := functions.Request(res, req, &request)

	if err != nil {
		return
	}

	user := new(models.User)
	err = user.FindByUsername(request.Username)
	if err != nil {
		functions.Response(res, req, 404, []byte("{}"))
		return
	}

	if user.Password != request.Password {
		functions.Response(res, req, 401, []byte("{}"))
		return
	}
	access, _ := jwt.GenerateToken("access", user.Username)
	refresh, _ := jwt.GenerateToken("refresh", user.Username)
	response := responses.SignInResponse{Access: access, Refresh: refresh}

	b, _ := json.MarshalIndent(response, "", "  ")

	functions.Response(res, req, 200, b)
	return
}

func (c AuthController) refresh(res http.ResponseWriter, req *http.Request) {
	claims := jwt.AuthRequire(res, req, "refresh")
	if claims == nil {
		return
	}

	access, _ := jwt.GenerateToken(claims.Identity, "access")
	response := responses.SignInResponse{Access: access}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()).Hours() <= 168 {
		response.Refresh, _ = jwt.GenerateToken(claims.Identity, "refresh")
	}

	b, _ := json.MarshalIndent(response, "", "  ")

	*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 200))
	res.WriteHeader(200)
	res.Write(b)
	return
}