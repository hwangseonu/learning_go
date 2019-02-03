package test

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/hwangseonu/goBackend/common/jwt"
	"github.com/hwangseonu/goBackend/users/controllers"
	"github.com/hwangseonu/goBackend/users/responses"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"os"
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

	res := responses.SignInResponse{}
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Error(err)
	} else {
		token, err := jwt.ParseWithClaims(res.Access, &jwt2.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			t.Error(err)
		}
		claims := token.Claims.(*jwt2.CustomClaims)

		if claims.Valid() != nil {
			t.Error("jwt is invalid")
		}

		token, err = jwt.ParseWithClaims(res.Refresh, &jwt2.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			t.Error(err)
		}
		claims = token.Claims.(*jwt2.CustomClaims)

		if claims.Valid() != nil {
			t.Error("jwt is invalid")
		}
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
