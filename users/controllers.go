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
	var request SignInRequest
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

	user := new(User)
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
	access, _ := GenerateToken("access", user.Username)
	refresh, _ := GenerateToken("refresh", user.Username)
	response := SignInResponse{access, refresh}

	b, _ := json.MarshalIndent(response, "", "  ")

	*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", 200))
	res.WriteHeader(200)
	res.Write(b)
	return
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