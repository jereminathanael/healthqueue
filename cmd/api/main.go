package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jereminathanael/healthqueue/internal/cache"
	"github.com/jereminathanael/healthqueue/internal/config"
	"github.com/jereminathanael/healthqueue/internal/database"
)

func main() {
	// Load config dari .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// koneksi ke db(postgreSQL)
	db, err := database.Connection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL connected")

	// koneksi ke redis
	cacheClient, err := cache.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer cacheClient.Close()
	log.Println("✅ Redis connected")

	// setup router
	r := chi.NewRouter()

	// Built-in middleware chi
	r.Use(middleware.Logger) // log setiap request
	r.Use(middleware.Recoverer) // recover dari panic
	r.Use(middleware.RequestID) // tambah request ID ke setiap request

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok, "service": "healthqueue"}`)
	})

	// API routes akan di tambah di tahap berikut nya
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"message": "HealthQueue API v1"}`)
		})
	})

	// start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server running on http://localhost/api/v1%s", addr)
  log.Printf("📋 Environment: %s", cfg.AppEnv)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}