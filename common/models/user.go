package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Username string
	Password string
	Nickname string
	Email    string
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
