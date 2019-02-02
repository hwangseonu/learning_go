package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	Id      int    `json:"id" bson:"_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Writer bson.ObjectId `json:"writer"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (post *Post) New(title, content string, writer *User) {
	post.Id = GetNextId("posts")
	post.Title = title
	post.Content = content
	post.Writer = writer.Id
	post.CreateAt = time.Now()
	post.UpdateAt = time.Now()
}

func (post *Post) FindById(id int) error {
	s := session.Clone()
	defer s.Close()

	var result Post
	if err := s.DB("backend").C("posts").FindId(id).One(&result); err != nil {
		return err
	}

	*post = result
	return nil
}

func (post *Post) Save() error {
	s := session.Clone()
	defer s.Close()

	_, err := s.DB("backend").C("posts").Upsert(bson.M{"_id": post.Id}, post)
	return err
}