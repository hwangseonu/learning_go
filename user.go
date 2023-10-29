package main

import (
	"errors"
	"net/http"
)

type User struct {
	Username string
	Password string `json:"-"`
	Nickname string
}

type UserAPI struct {
	API
}

func NewUserAPI() *UserAPI {
	api := &UserAPI{}
	api.HandleFunc(http.MethodPost, "/users", api.CreateUser)

	return api
}

func (api *UserAPI) CreateUser(writer http.ResponseWriter, request *http.Request) {
	body := struct {
		Username string
		Password string
		Nickname string
	}{}

	if err := api.BodyMapping(request, &body); err != nil {
		api.RespondError(writer, http.StatusInternalServerError, err)
		return
	}

	passwordHash := hashString([]byte(body.Password))
	user := User{body.Username, passwordHash, body.Nickname}

	if _, ok := store.Users[user.Username]; ok {
		api.RespondError(writer, http.StatusConflict, errors.New("users information cannot be saved. already exists `Username`"))
	}

	store.Users[user.Username] = user
	api.RespondJSON(writer, http.StatusOK, user)
}
