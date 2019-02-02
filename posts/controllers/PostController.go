package controllers

import (
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/functions"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/common/models"
	"github.com/hwangseonu/goBackend/posts/requests"
	"github.com/hwangseonu/goBackend/posts/responses"
	userRes"github.com/hwangseonu/goBackend/users/responses"
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
	} else if regexp.MustCompile(`^/posts/\d+$`).Match(path) && req.Method == "DELETE" {
		id, err := strconv.Atoi(strings.TrimPrefix(req.URL.Path, "/posts/"))
		if err != nil {
			functions.Response(res, req, 400, []byte(`{"message": "post is is integer"}`))
			return
		}
		c.deletePost(res, req, id)
	} else {
		functions.Response(res, req, 404, []byte(`{"message": "404 page not found"}`))
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
	post := new(models.Post)
	if err := post.FindById(id); err != nil {
		functions.Response(res, req, 404, []byte(`{"message": "cannot find post by id"}`))
		return
	}

	response := responses.PostResponses{
		Id: post.Id,
		Title: post.Title,
		Content: post.Content,
		CreateAt: post.CreateAt,
		UpdateAt: post.UpdateAt,
	}

	user := new(models.User)
	if err := user.FindById(post.Writer); err != nil {
		response.Writer = userRes.GetUserResponse{}
	} else {
		response.Writer = userRes.GetUserResponse{Username: user.Username, Nickname: user.Nickname, Email: user.Email}
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	functions.Response(res, req, 200, b)
}

func (c PostController) deletePost(res http.ResponseWriter, req *http.Request, id int) {
	claims := jwt.AuthRequire(res, req, "access")
	if claims == nil {
		return
	}
	user := new(models.User)
	user.FindByUsername(claims.Identity)

	post := new(models.Post)
	if err := post.FindById(id); err != nil {
		functions.Response(res, req, 404, []byte(`{"message": "cannot find post by id"}`))
		return
	}

	if post.Writer.Hex() != user.Id.Hex() {
		functions.Response(res, req, 403, []byte(`{"message": "this post is not your own"}`))
		return
	}

	if err := post.Delete(); err != nil {
		functions.Response(res, req, 404, []byte(`{"message": "`+err.Error()+`"}`))
		return
	}
	functions.Response(res, req, 200, []byte(`{}`))
}