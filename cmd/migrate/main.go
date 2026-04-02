package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/example/ordersvc/internal/config"
	_ "github.com/lib/pq"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`,
	`CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(100) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price_cents INTEGER NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0,
		category VARCHAR(100),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`,
	`CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		total_cents INTEGER NOT NULL DEFAULT 0,
		shipping_address TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`,
	`CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INTEGER REFERENCES orders(id),
		product_id INTEGER REFERENCES products(id),
		quantity INTEGER NOT NULL,
		price_cents INTEGER NOT NULL
	)`,
}

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	switch direction {
	case "up":
		runMigrations(db)
	case "down":
		log.Println("down migrations not implemented yet")
	default:
		log.Fatalf("unknown migration direction: %s", direction)
	}
}

func runMigrations(db *sql.DB) {
	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Fatalf("migration %d failed: %v", i+1, err)
		}
		fmt.Printf("migration %d applied successfully\n", i+1)
	}
	fmt.Println("all migrations applied")
}
