package main

import (
	"github.com/hwangseonu/goBackend/common"
	"github.com/hwangseonu/goBackend/users"
	"net/http"
)

const port = ":5000"

func main() {
	before := common.ChainBeforeMiddlewares(common.LoggerMiddleware,)
	after := common.ChainAfterMiddlewares(common.LoggerMiddleware)

	http.HandleFunc("/auth", before(after(new(users.AuthController))))
	http.Handle("/users", before(after(new(users.UserController))))

	println("server is running on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}