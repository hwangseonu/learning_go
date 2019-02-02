package controllers

import (
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/functions"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/posts/requests"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type PostController struct {
	http.Handler
}

func (c *PostController) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := []byte(req.URL.Path)

	if regexp.MustCompile("^/posts$").Match(path) && req.Method == "POST" {
		c.createPost(res, req)
	} else if regexp.MustCompile(`^/posts/\d+$`).Match(path) && req.Method == "GET" {
		id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/posts/"))
		if err != nil {
			functions.Response(res, req, 400, []byte(`{"message": "post is is integer"}`))
			return
		}
		c.getPost(res, req, id)
	}
}

func (c PostController) createPost(res http.ResponseWriter, req *http.Request) {
	claims := jwt.AuthRequire(res, req, "access")
	if claims == nil {
		return
	}

	user := new(models.User)
	user.FindByUsername(claims.Identity)

	var request requests.CreatePostRequest
	err := functions.Request(res, req, &request)

	if err != nil {
		return
	}

	post := new(models.Post)
	post.New(request.Title, request.Content, user)
	err = post.Save()

	if err != nil {
		functions.Response(res, req, 500, []byte(`{"message": "`+err.Error()+`"}`))
		return
	}

	b, _ := json.MarshalIndent(post, "", "  ")
	functions.Response(res, req, 201, b)
	return
}

func (c PostController) getPost(res http.ResponseWriter, req *http.Request, id int) {
	functions.Response(res, req, 200, []byte(`{"message": "`+strconv.Itoa(id)+`"}`))
}