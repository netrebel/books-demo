package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
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
	assert.Equal(t, []byte("{\"id\":0,\"isbn\":\"1234\",\"title\":\"Life\",\"author\":{\"first_name\":\"Angel\",\"last_name\":\"Reyes\"}}\n"), body, "Should return created Book")
}

// TestUpdateBookSuccess will only run after POST
func TestUpdateBookSuccess(t *testing.T) {
	r := strings.NewReader("{\"isbn\":\"1234\",\"title\":\"Life\",\"author\":{\"first_name\":\"Miguel\",\"last_name\":\"Reyes\"}}")
	request, _ := http.NewRequest("PUT", "/books/0", r)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "200 response is expected")

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("HERE: ", string(body))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, []byte("{\"id\":0,\"isbn\":\"1234\",\"title\":\"Life\",\"author\":{\"first_name\":\"Miguel\",\"last_name\":\"Reyes\"}}\n"), body, "Should return created Book")
}

// TestDeleteBookByIdOK will only run after POST
func TestDeleteBookByIdOK(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/books/0", nil)
	resp := httptest.NewRecorder()
	Router().ServeHTTP(resp, request)
	assert.Equal(t, http.StatusOK, resp.Code, "404 response is expected")
}
