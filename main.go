package main

import (
	"github.com/hwangseonu/goBackend/server/controller"
	"net/http"
)

const port = ":5000"

func main() {
	http.Handle("/auth", new(controller.AuthController))
	println("server is running on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}