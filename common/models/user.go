package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
