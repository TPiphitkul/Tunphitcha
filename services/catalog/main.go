package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("catalog:ok"))
	})
	r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"items":["apple","banana"]}`))
	})
	log.Println("catalog service on :8082")
	http.ListenAndServe(":8082", r)
}
