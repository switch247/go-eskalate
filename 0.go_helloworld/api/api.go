package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// In-memory storage
var (
	books  = []Book{}
	nextID = 1
	mu     sync.Mutex
)

// Handler to add a new book
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the Function.
	var body = json.NewDecoder(r.Body)
	NewBook := Book{ID: nextID, Title: body.Title, Author: body.Author, Year: body.Year}
	books = append(books, NewBook)
	nextID++
	return json.Encoder(NewBook.ID)
}

func RunServer() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		// Route requests to the appropriate handler based on the HTTP method
		if r.Method == http.MethodPost {
			http.Post(w, addBookHandler(w, r), http.StatusCreated)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
