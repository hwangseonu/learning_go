package test

import (
	"encoding/json"
	"github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/posts/controllers"
	"github.com/hwangseonu/goBackend/posts/responses"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func CreatePost(t *testing.T) int {
	access, err := jwt.GenerateToken("access", "test")

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/posts", strings.NewReader(`{"title": "test", "content": "test"}`))
	req.Header.Add("Authorization", "Bearer "+access)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	res := responses.PostResponse{}
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	} else {
		return res.Id
	}
	return 0
}

func DropPosts(t *testing.T) {
	if s, err := mgo.Dial("mongodb://localhost:27017"); err != nil {
		t.Fatal(err)
	} else {
		if err := s.DB("backend").C("posts").DropCollection(); err != nil {
			t.Fatal(err)
		}
		s.Close()
	}
}

func TestCreatePost(t *testing.T) {
	SignUp("test", t)
	CreatePost(t)
	DropPosts(t)
	RemoveTestUser(t)
}

func TestGetPost(t *testing.T) {
	SignUp("test", t)

	//Check not found
	req, err := http.NewRequest("GET", "/posts/"+strconv.Itoa(0), nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	id := CreatePost(t)

	//Check ok
	req, err = http.NewRequest("GET", "/posts/"+strconv.Itoa(id), nil)

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	DropPosts(t)
	RemoveTestUser(t)
}

func TestUpdatePost(t *testing.T) {
	SignUp("test", t)
	SignUp("test1", t)

	access, err := jwt.GenerateToken("access", "test")

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "/posts/0", strings.NewReader(`{"title": "test1234", "content": "test"}`))
	req.Header.Add("Authorization", "Bearer "+access)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	id := CreatePost(t)

	req, err = http.NewRequest("PATCH", "/posts/"+strconv.Itoa(id), strings.NewReader(`{"title": "test1234", "content": "test"}`))
	req.Header.Add("Authorization", "Bearer "+access)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	res := responses.PostResponse{}
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	} else {
		if res.Title != "test1234" {
			t.Errorf("The title of the post had to be changed to test1234. post title is %s", res.Title)
		}
	}

	if access, err = jwt.GenerateToken("access", "test1"); err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("PATCH", "/posts/"+strconv.Itoa(id), strings.NewReader(`{"title": "test1234", "content": "test"}`))
	req.Header.Add("Authorization", "Bearer "+access)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}

	DropPosts(t)
	RemoveTestUser(t)
}

func TestDeletePost(t *testing.T) {
	SignUp("test", t)
	SignUp("test1", t)

	access, err := jwt.GenerateToken("access", "test1")

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/posts/0", nil)
	req.Header.Add("Authorization", "Bearer "+access)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	id := CreatePost(t)

	req, err = http.NewRequest("DELETE", "/posts/"+strconv.Itoa(id), nil)
	req.Header.Add("Authorization", "Bearer "+access)

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}

	access, err = jwt.GenerateToken("access", "test")

	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("DELETE", "/posts/"+strconv.Itoa(id), nil)
	req.Header.Add("Authorization", "Bearer "+access)

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = new(controllers.PostController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	DropPosts(t)
	RemoveTestUser(t)
}