package models

import (
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

func GetNextId(doc string) int {
	s := session.Copy()
	defer s.Close()
	DB := s.DB("backend")
	var counter map[string]interface{}

	if err := DB.C("auto_increment").Find(bson.M{"document": doc}).One(&counter); err != nil || counter == nil {
		DB.C("auto_increment").Insert(map[string]interface{}{
			"count": 1,
			"document": doc,
		})
		return 0
	} else {
		id := counter["count"].(int)
		counter["count"] = id + 1
		DB.C("auto_increment").Update(bson.M{"document": doc}, counter)
		return id
	}
}
