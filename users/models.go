package users

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

var session *mgo.Session

func init() {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	session = s
}

type User struct {
	Username string
	Password string
	Nickname string
	Email string
}

func (u *User) FindByUsername(username string) error {
	s := session.Clone()
	defer s.Close()
	var result User
	err := s.DB("backend").C("users").Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return err
	}
	*u = result
	return nil
}

func (u User) Save() error {
	s := session.Clone()
	defer s.Close()
	err := (&User{}).FindByUsername(u.Username)

	if err == nil {
		return fmt.Errorf("user already exists")
	}

	err = s.DB("backend").C("users").Insert(u)
	return err
}

type CustomClaims struct {
	jwt.StandardClaims
	Identity string `json:"identity"`
}

var secret = os.Getenv("jwt-secret")

func (c CustomClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	user := new(User)
	err := user.FindByUsername(c.Identity)

	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("could not find user")
	}
	return nil
}

func GenerateToken(t, username string) (string, error) {
	var expire int64

	if t == "access" {
		expire = time.Now().Add(time.Hour).Unix()
	} else {
		expire = time.Now().AddDate(0, 1, 0).Unix()
	}

	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: expire,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "",
		Subject: t,
		NotBefore: time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS512, CustomClaims{claims, username}).SignedString([]byte(secret))
}
