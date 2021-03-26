package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/netrebel/books-demo/handlers"
	"github.com/netrebel/books-demo/promutil"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(promutil.PrometheusMiddleware)
	// Serving static files
	router.Handle("/", nil).Handler(http.FileServer(http.Dir("./static/")))
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/books", handlers.AddBook).Methods("POST")
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

	// Prometheus endpoint
	router.Path("/prometheus").Handler(promhttp.Handler())

	log.Println("Listening on http://localhost:9000/")
	http.ListenAndServe(":9000", router)
}
