package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user:ok"))
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok","user":"demo"}`))
	})
	log.Println("user service on :8081")
	http.ListenAndServe(":8081", r)
}
