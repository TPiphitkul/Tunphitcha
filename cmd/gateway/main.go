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

// adaptiveMiddleware ประมวลผลคำขอพร้อมส่งข้อมูลให้ Risk Analyzer
func adaptiveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ดึงข้อมูลเมตาจาก request
		meta := profiler.Extract(r)

		// ใช้ ML Risk Analyzer ประเมินความเสี่ยง
		score := risk.MLScore(meta)
		level := risk.Level(score)
		dec := policy.Decide(level)

		// เพิ่ม header เพื่อ debug
		w.Header().Set("X-Risk-Score", fmt.Sprint(score))
		w.Header().Set("X-Risk-Level", level)

		// ใช้นโยบาย enforcement ตามระดับ risk
		enf := enforcement.Apply(dec.RateLimit)
		enf(next).ServeHTTP(w, r)

		// เก็บ log metadata เพื่อนำไปเทรนภายหลัง
		meta.RiskScore = score
		profiler.Export(meta)
	})
}

func main() {
	// เริ่มเชื่อมต่อ Redis (สำหรับ rate-limit และ profiling)
	profiler.InitRedis("redis:6379")
	enforcement.InitRedis("redis:6379")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Endpoint ตรวจสอบสถานะระบบ
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("gw:ok"))
	})

	// Static baseline routes (ไม่ผ่าน ML)
	r.Mount("/static", BackendRoutes())

	// Adaptive routes (ผ่าน ML Risk Analyzer)
	r.Group(func(r chi.Router) {
		r.Use(adaptiveMiddleware)
		r.Mount("/api", BackendRoutes())
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown
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
	log.Println("Gateway shutdown complete")
}
