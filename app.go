package main

import (
	. "book-collections-restapi/config"
	. "book-collections-restapi/dao"
	"book-collections-restapi/dto/request"
	"book-collections-restapi/helpers"
	. "book-collections-restapi/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var dao BookDao

func AllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dao.FindAll()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, books)
}

func FindBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := dao.FindOne(params["id"])

	if book.IsEmpty() {
		helpers.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bookRequest := request.BookRequest{}

	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := bookRequest.Validate(w, r)
	if err {
		return
	}

	var book Book

	copier.Copy(&book, &bookRequest)

	book.ID = bson.NewObjectId()
	if err := dao.Insert(book); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusCreated, book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bookRequest := request.BookRequest{}

	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := bookRequest.Validate(w, r)
	if err {
		return
	}

	var book Book

	copier.Copy(&book, &bookRequest)

	if err := dao.Update(book); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := dao.FindOne(params["id"])

	if book.IsEmpty() {
		helpers.RespondWithError(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := dao.Delete(book); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, map[string]string{"success": "the book has been deleted"})
}

func init() {
	config := Config{}
	config.Read()

	dao = BookDao{}
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/books", AllBooks).Methods("GET")
	router.HandleFunc("/books", CreateBook).Methods("POST")
	router.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", FindBook).Methods("GET")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
