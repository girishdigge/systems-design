package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	TargetUsers    = 100000
	TargetProducts = 10000
	TargetOrders   = 1000000
	TargetEvents   = 5000000
	BatchSize      = 50000
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, falling back to system environment variables")
	}
	os.Unsetenv("PGPASSWORD")
	os.Unsetenv("PGUSER")
	os.Unsetenv("PGDATABASE")
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Starting mass ingestion sequence...")

	// 1. Seed Products
	productIDs := make([]uuid.UUID, TargetProducts)
	var productRows [][]any
	for i := range TargetProducts {
		productIDs[i] = uuid.New()
		productRows = append(productRows, []any{
			productIDs[i],
			fmt.Sprintf("Product Premium Specifications #%d", i),
			10.0 + rand.Float64()*990.0,
		})
	}
	copyRows(ctx, pool, "products", []string{"id", "name", "price"}, productRows)

	// 2. Seed Users
	userIDs := make([]uuid.UUID, TargetUsers)
	var userRows [][]any
	for i := range TargetUsers {
		userIDs[i] = uuid.New()
		userRows = append(userRows, []any{
			userIDs[i],
			fmt.Sprintf("Engineer Candidate %d", i),
			fmt.Sprintf("backend.engine.%d@scalelab.internal", i),
			time.Now().Add(-time.Duration(rand.Intn(10000)) * time.Hour),
		})

		if len(userRows) >= BatchSize || i == TargetUsers-1 {
			copyRows(ctx, pool, "users", []string{"id", "name", "email", "created_at"}, userRows)
			userRows = nil
		}
	}

	// 3. Seed Orders & Order Items
	var orderRows [][]any
	var itemRows [][]any

	for i := range TargetOrders {
		orderID := uuid.New()
		userID := userIDs[rand.Intn(len(userIDs))]
		orderTime := time.Now().Add(-time.Duration(rand.Intn(5000)) * time.Hour)

		orderRows = append(orderRows, []any{
			orderID,
			userID,
			0.0, // Updated later or managed at analytical runtime
			orderTime,
		})

		// 1 to 3 items per order
		itemsCount := rand.Intn(3) + 1
		for range itemsCount {
			itemRows = append(itemRows, []any{
				uuid.New(),
				orderID,
				productIDs[rand.Intn(len(productIDs))],
				rand.Intn(5) + 1,
			})
		}

		if len(orderRows) >= BatchSize || i == TargetOrders-1 {
			copyRows(ctx, pool, "orders", []string{"id", "user_id", "total", "created_at"}, orderRows)
			copyRows(ctx, pool, "order_items", []string{"id", "order_id", "product_id", "quantity"}, itemRows)
			orderRows = nil
			itemRows = nil
			log.Printf("Ingested %d / %d Orders...", i+1, TargetOrders)
		}
	}

	// 4. Seed System Events
	eventTypes := []string{"USER_LOGIN", "ITEM_VIEWED", "CHECKOUT_ABANDONED", "PAYMENT_ATTEMPTED", "SESSION_EXPIRED"}
	var eventRows [][]any
	for i := range TargetEvents {
		eventRows = append(eventRows, []any{
			userIDs[rand.Intn(len(userIDs))],
			eventTypes[rand.Intn(len(eventTypes))],
			time.Now().Add(-time.Duration(rand.Intn(8000)) * time.Hour),
		})

		if len(eventRows) >= BatchSize || i == TargetEvents-1 {
			copyRows(ctx, pool, "events", []string{"user_id", "event_type", "created_at"}, eventRows)
			eventRows = nil
		}
	}
	log.Println("Database Lab population complete.")
}

func copyRows(ctx context.Context, pool *pgxpool.Pool, tableName string, columns []string, rows [][]any) {
	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatalf("Fatal CopyFrom bulk insertion error into %s: %v", tableName, err)
	}
}
