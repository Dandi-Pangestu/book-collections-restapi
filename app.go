package main

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	. "book-collections-restapi/dao"
	. "book-collections-restapi/models"
	. "book-collections-restapi/config"
)

var dao BookDao

func AllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, books)
}

func FindBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := dao.FindOne(params["id"])

	if book.IsEmpty() {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	book.ID = bson.NewObjectId()
	if err := dao.Insert(book); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dao.Update(book); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := dao.FindOne(params["id"])

	if book.IsEmpty() {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := dao.Delete(book); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{"success": "the book has been deleted"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
