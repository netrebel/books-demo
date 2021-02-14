package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Welcome!</h1>")
}

func main() {
	mux := mux.NewRouter().StrictSlash(true)
	mux.HandleFunc("/", home)
	fmt.Println("Listening on http://localhost:8080/")
	http.ListenAndServe(":8080", mux)
}
