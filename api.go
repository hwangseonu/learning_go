package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type route struct {
	method  string
	pattern string
	handler http.Handler
}

func (r *route) match(method, path string) bool {
	if r.method != method {
		return false
	}

	splitPattern := strings.Split(r.pattern, "/")
	splitPath := strings.Split(path, "/")

	if len(splitPattern) != len(splitPath) {
		return false
	}

	for i, p1 := range splitPattern {
		p2 := splitPath[i]

		if !strings.HasPrefix(":", p1) && p1 != p2 {
			return false
		}
	}

	return true
}

type API struct {
	routes []route
}

func (api *API) Handle(method string, pattern string, handler http.Handler) {
	r := route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	api.routes = append(api.routes, r)
}

func (api *API) HandleFunc(method string, pattern string, handler http.HandlerFunc) {
	api.Handle(method, pattern, handler)
}

func (api *API) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path

	for _, r := range api.routes {
		if r.match(method, path) {
			r.handler.ServeHTTP(writer, request)
			return
		}
	}

	api.RespondError(writer, http.StatusNotFound, errors.New("404 page not found"))
	return
}

func (api *API) BodyMapping(request *http.Request, v interface{}) error {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	return err
}

func (api *API) RespondError(writer http.ResponseWriter, status int, err error) {
	log.Printf("error: %v\n", err)
	body := map[string]string{"message": err.Error()}
	api.RespondJSON(writer, status, body)
}

func (api *API) RespondJSON(writer http.ResponseWriter, status int, body interface{}) {
	res, err := json.Marshal(body)
	if err != nil {
		api.RespondError(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	_, err = writer.Write(res)
	if err != nil {
		api.RespondError(writer, http.StatusInternalServerError, err)
		return
	}
}
