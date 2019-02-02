package main

import (
	"github.com/hwangseonu/goBackend/common/middlewares"
	posts "github.com/hwangseonu/goBackend/posts/controllers"
	users "github.com/hwangseonu/goBackend/users/controllers"
	"net/http"
)

const port = ":5000"

func main() {
	before := middlewares.ChainBeforeMiddlewares()
	after := middlewares.ChainAfterMiddlewares(middlewares.LoggerMiddleware)

	http.HandleFunc("/auth", before(after(new(users.AuthController))))
	http.HandleFunc("/auth/refresh", before(after(new(users.AuthController))))
	http.HandleFunc("/users", before(after(new(users.UserController))))
	http.Handle("/posts", before(after(new(posts.PostController))))

	println("server is running on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}