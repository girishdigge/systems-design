package main

import (
	"database-scaling-lab/internal/db"
	"database-scaling-lab/internal/handlers"
	"database-scaling-lab/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, falling back to system environment variables")
	}
	os.Unsetenv("PGPASSWORD")
	os.Unsetenv("PGUSER")
	os.Unsetenv("PGDATABASE")

	dsn := fmt.Sprintf("postgres://%s:%s@postgres-primary:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	pool, err := db.NewPool(dsn)
	if err != nil {
		log.Fatalf("Critical pool initialization failure: %v", err)
	}
	defer pool.Close()

	repo := &repository.PostgresRepository{Pool: pool}
	server := &handlers.Server{Repo: repo}

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	// Expose dedicated internal scraping endpoint for Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Database Lab API Operational on cluster port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server aborted: %v", err)
	}
}
