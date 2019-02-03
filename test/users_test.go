package test

import (
	"github.com/hwangseonu/goBackend/users/controllers"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignUp(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", strings.NewReader(`{"username": "test", "password": "test1234", "nickname": "test", "email": "test@test"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.UserController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	req, err = http.NewRequest("POST", "/users", strings.NewReader(`{"username": "test", "password": "test1234", "nickname": "test", "email": "test@test"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	if s, err := mgo.Dial("mongodb://localhost:27017"); err != nil {
		t.Error(err)
	} else {
		if err = s.DB("backend").C("users").Remove(bson.M{"username": "test"}); err != nil {
			t.Error(err)
			s.Close()
		}
	}
}

func TestSignIn(t *testing.T) {
	req, err := http.NewRequest("POST", "/users", strings.NewReader(`{"username": "test", "password": "test1234", "nickname": "test", "email": "test@test"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.Handler(new(controllers.UserController))

	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("POST", "/auth", strings.NewReader(`{"username": "test", "password": "test1234"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.Handler(new(controllers.AuthController))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if s, err := mgo.Dial("mongodb://localhost:27017"); err != nil {
		t.Error(err)
	} else {
		if err = s.DB("backend").C("users").Remove(bson.M{"username": "test"}); err != nil {
			t.Error(err)
			s.Close()
		}
	}
}
