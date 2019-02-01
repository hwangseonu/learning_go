package main

import (
	"github.com/hwangseonu/goBackend/users"
	"net/http"
)

const port = ":5000"

func main() {
	http.Handle("/auth", new(users.AuthController))
	http.Handle("/users", new(users.UserController))
	println("server is running on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}