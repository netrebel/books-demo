package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/netrebel/books-demo/redisdb"
	uuid "github.com/satori/go.uuid"
)

var redisDb = redisdb.NewDatabase()
var books []Book

func init() {
	_, err := redisDb.Client.Ping(redisdb.Ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

// Book struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Used for logging purposes
func (a *Author) String() string {
	return fmt.Sprintf("{ FirstName : %s, LastName: %s}", a.FirstName, a.LastName)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if books == nil {
		books = []Book{}
	}
	log.Printf("%+v\n", books)
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find book Id
	vars := mux.Vars(r)
	bookID := vars["id"]

	val, err := redisDb.Client.Get(redisdb.Ctx, bookID).Result()
	if err == redis.Nil {
		log.Println("Cache miss for id: ", vars["id"])
		for _, book := range books {
			if book.ID == vars["id"] {
				log.Printf("Found:  %+v\n", book)
				json.NewEncoder(w).Encode(book)
				return
			}
		}
	} else {
		log.Printf("Redis hit: %v\n", val)
		var b Book
		json.Unmarshal([]byte(val), &b)
		json.NewEncoder(w).Encode(b)
		return
	}
	log.Println("Not found!")
	http.Error(w, "404 not found", http.StatusNotFound)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	log.Printf("POST payload: %+v\n", book)

	book.ID = uuid.NewV4().String()

	// Append to list
	books = append(books, book)

	// Save to redis
	m, err := json.Marshal(book)
	if err != nil {
		log.Fatalln(err)
	}
	redisDb.Client.Set(redisdb.Ctx, book.ID, string(m), 24*time.Hour)
	log.Printf("%+v\n", books)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var updatedBook Book
	for index, book := range books {
		if book.ID == vars["id"] {
			_ = json.NewDecoder(r.Body).Decode(&updatedBook)
			log.Printf("PUT payload:  %+v\n", updatedBook)

			updatedBook.ID = book.ID
			books[index] = updatedBook
			json.NewEncoder(w).Encode(books[index])
			return
		}
		json.NewEncoder(w).Encode(&Book{})
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	var deletedBook Book
	for index, book := range books {
		if book.ID == vars["id"] {
			deletedBook = books[index]

			books = append(books[:index], books[index+1:]...)
			_ = json.NewEncoder(w).Encode(&deletedBook)
			return
		}
	}
	http.Error(w, "404 not found", http.StatusNotFound)
}
