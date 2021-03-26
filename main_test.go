package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/netrebel/books-demo/handlers"
	"github.com/stretchr/testify/assert"
)

var bookID string

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")
	router.HandleFunc("/books", handlers.AddBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	return router
}

func TestGetBooksEmpty(t *testing.T) {
	request, _ := http.NewRequest("GET", "/books", nil)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "OK response is expected")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []byte("[]\n"), body, "Should return Empty list []")
}

func TestGetBookByIdNotFound(t *testing.T) {
	request, _ := http.NewRequest("GET", "/books/{id}", nil)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusNotFound, resp.Code, "404 response is expected")
}
func TestDeleteBookByIdNotFound(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/books/0", nil)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusNotFound, resp.Code, "404 response is expected")
}
func TestAddBookSuccess(t *testing.T) {
	r := strings.NewReader("{\"isbn\":\"1234\",\"title\":\"Life\",\"author\":{\"first_name\":\"Angel\",\"last_name\":\"Reyes\"}}")
	request, _ := http.NewRequest("POST", "/books", r)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "200 response is expected")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var insertedBook handlers.Book
	err = json.Unmarshal(body, &insertedBook)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
	}
	bookID = insertedBook.ID
	assert.Equal(t, "1234", insertedBook.Isbn, "Isbn did not match")
	assert.Equal(t, "Life", insertedBook.Title, "Title did not match")
	assert.Equal(t, "Angel", insertedBook.Author.FirstName, "FirstName did not match")
	assert.Equal(t, "Reyes", insertedBook.Author.LastName, "LastName did not match")
}

// TestUpdateBookSuccess will only run after POST
func TestUpdateBookSuccess(t *testing.T) {
	r := strings.NewReader("{\"isbn\":\"1234\",\"title\":\"Life\",\"author\":{\"first_name\":\"Miguel\",\"last_name\":\"Reyes\"}}")
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/books/%v", bookID), r)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "200 response is expected")

	var updatedBook handlers.Book
	_ = json.NewDecoder(resp.Body).Decode(&updatedBook)

	assert.Equal(t, "1234", updatedBook.Isbn, "Isbn did not match")
	assert.Equal(t, "Life", updatedBook.Title, "Title did not match")
	assert.Equal(t, "Miguel", updatedBook.Author.FirstName, "FirstName did not match")
	assert.Equal(t, "Reyes", updatedBook.Author.LastName, "LastName did not match")
}

// TestDeleteBookByIdOK will only run after POST
func TestDeleteBookByIdOK(t *testing.T) {
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/books/%v", bookID), nil)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "200 response is expected")
}
