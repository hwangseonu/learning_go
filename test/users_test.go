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

func SignUp(t *testing.T) {
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
}

func RemoveTestUser(t *testing.T) {
	if s, err := mgo.Dial("mongodb://localhost:27017"); err != nil {
		t.Fatal(err)
	} else {
		if err = s.DB("backend").C("users").Remove(bson.M{"username": "test"}); err != nil {
			t.Fatal(err)
			s.Close()
		}
	}
}

func JwtCheck(tokenString string, t *testing.T) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt2.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Error(err)
	}
	claims := token.Claims.(*jwt2.CustomClaims)

	if claims.Valid() != nil {
		t.Error("jwt is invalid")
	}
}

func TestSignUp(t *testing.T) {
	SignUp(t)

	//Test user conflict
	req, err := http.NewRequest("POST", "/users", strings.NewReader(`{"username": "test", "password": "test1234", "nickname": "test", "email": "test@test"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.UserController)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	RemoveTestUser(t)
}

func TestSignIn(t *testing.T) {
	SignUp(t)

	req, err := http.NewRequest("POST", "/auth", strings.NewReader(`{"username": "test", "password": "test1234"}`))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.Handler(new(controllers.AuthController))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	res := responses.SignInResponse{}
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	} else {
		JwtCheck(res.Access, t)
		JwtCheck(res.Refresh, t)
	}

	RemoveTestUser(t)
}

func TestGetUser(t *testing.T) {
	access, err := jwt2.GenerateToken("access", "test")

	if err != nil {
		t.Error(err)
	}

	//Check unprocessable entity
	req, err := http.NewRequest("GET", "/users", nil)
	req.Header.Add("Authorization", "Bearer " + access)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := new(controllers.UserController)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	SignUp(t)

	//Check ok
	req, err = http.NewRequest("GET", "/users", nil)
	req.Header.Add("Authorization", "Bearer " + access)

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	RemoveTestUser(t)
}