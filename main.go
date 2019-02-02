package main

import (
	"github.com/hwangseonu/goBackend/common/middlewares"
	"github.com/hwangseonu/goBackend/users/controllers"
	"net/http"
)

const port = ":5000"

func main() {
	before := middlewares.ChainBeforeMiddlewares()
	after := middlewares.ChainAfterMiddlewares(middlewares.LoggerMiddleware)

	http.HandleFunc("/auth", before(after(new(controllers.AuthController))))
	http.Handle("/users", before(after(new(controllers.UserController))))

	println("server is running on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}