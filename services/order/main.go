package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("order:ok"))
	})
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"order":"created"}`))
	})
	log.Println("order service on :8083")
	http.ListenAndServe(":8083", r)
}
