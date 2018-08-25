package models

import "gopkg.in/mgo.v2/bson"

type Book struct {
	ID bson.ObjectId `bson:"_id" json:"id"`
	Title string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
	Publisher string `bson:"publisher" json:"publisher"`
	Description string `bson:"description" json:"description"`
}

func (book *Book) IsEmpty() bool {
	return *book == Book{}
}