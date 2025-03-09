package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lithammer/shortuuid/v4"
)

type Mapper struct {
	Mapping map[string]string
	Lock    sync.Mutex // Prevent deadlock
}

var urlMapper Mapper

func init() {
	urlMapper = Mapper{
		Mapping: make(map[string]string),
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Post("/shorten", createShortURLHandler)
	r.Get("/short/{key}", redirectHandler)
	http.ListenAndServe(":3000", r)
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // TODO: Parse JSON
	originalURL := r.Form.Get("url")

	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	key := shortuuid.New()

	insertMapping(key, originalURL)

	log.Println("Created mapping from", key, "to", originalURL)

	w.WriteHeader(http.StatusCreated)

	result := fmt.Sprintf("http://localhost:3000/short/%s", key)
	w.Write([]byte(result))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	originalURL, exists := readMapping(key)

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL not found"))
		return
	}

	// It doesn't work without a full URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// insertMapping adds a new mapping from a key to an original URL in the URL mapper.
// It acquires a lock on the URL mapper to ensure thread-safety during the insertion.
// The lock is released after the insertion is complete.
//
// This function is not concurrency-safe with respect to other operations on the URL mapper
// outside of the provided locking mechanism.
func insertMapping(key string, originalUrl string) {
	urlMapper.Lock.Lock()
	defer urlMapper.Lock.Unlock()

	urlMapper.Mapping[key] = originalUrl
}

func readMapping(key string) (string, bool) {
	urlMapper.Lock.Lock()
	defer urlMapper.Lock.Unlock()

	originalUrl, exists := urlMapper.Mapping[key]
	log.Println("para la key", key, "la url es", originalUrl)
	return originalUrl, exists
}
