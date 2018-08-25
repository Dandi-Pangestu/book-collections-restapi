package dao

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "book-collections-restapi/models"
)

type BookDao struct {
	Server string
	Database string
}

const COLLECTION string = "books"

var db *mgo.Database

func (dao *BookDao) Connect() {
	session, err := mgo.Dial(dao.Server)

	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(dao.Database)
}

func (dao *BookDao) FindAll() ([]Book, error) {
	var books []Book

	err := db.C(COLLECTION).Find(bson.M{}).All(&books)

	return books, err
}

func (dao *BookDao) FindOne(id string) (Book, error) {
	var book Book

	if !bson.IsObjectIdHex(id) {
		return book, nil
	}

	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&book)

	return book, err
}

func (dao *BookDao) Insert(book Book) error {
	err := db.C(COLLECTION).Insert(&book)
	return err
}

func (dao *BookDao) Delete(book Book) error {
	err := db.C(COLLECTION).Remove(&book)
	return err
}

func (dao *BookDao) Update(book Book) error {
	err := db.C(COLLECTION).UpdateId(book.ID, &book)
	return err
}