package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/example/go-adaptive-gw/internal/enforcement"
	"github.com/example/go-adaptive-gw/internal/policy"
	"github.com/example/go-adaptive-gw/internal/profiler"
	"github.com/example/go-adaptive-gw/internal/risk"
)



func BackendRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/catalog/list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"items":["apple","banana","cookie"]}`))
	})
	r.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	})
	r.Post("/order/create", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"order":"ok"}`))
	})
	return r
}

func adaptiveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract metadata
	meta := profiler.Extract(r)

	// NOTE: In this starter, ReqPerMin is not tracked; production should use Redis.
	meta.ReqPerMin = 1

	score := risk.Score(meta)
	level := risk.Level(score)
	dec := policy.Decide(level)

	w.Header().Set("X-Risk-Score", fmt.Sprint(score))
	w.Header().Set("X-Risk-Level", level)

	// Apply enforcement (wrap next)
	enf := enforcement.Apply(dec.RateLimit)
	enf(next).ServeHTTP(w, r)
	})
}

func main() {
	import "github.com/example/go-adaptive-gw/internal/enforcement"
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic health endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("gw:ok"))
	})

	// Static baseline mount
	r.Mount("/static", BackendRoutes())

	// Adaptive mount
	r.Group(func(r chi.Router) {
		r.Use(adaptiveMiddleware)
		r.Mount("/api", BackendRoutes())
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Gateway listening on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("Shutdown")
}
