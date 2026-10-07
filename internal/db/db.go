package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func EnsureSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(100) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		brand VARCHAR(100) NOT NULL DEFAULT 'Amazing',
		type VARCHAR(50) NOT NULL DEFAULT 'sneaker',
		category VARCHAR(100) NOT NULL,
		price NUMERIC(12, 2) NOT NULL,
		old_price NUMERIC(12, 2),
		rating NUMERIC(3, 1) DEFAULT 4.5,
		reviews INT DEFAULT 0,
		description TEXT,
		tags JSONB DEFAULT '[]'::jsonb,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS product_images (
		id SERIAL PRIMARY KEY,
		product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		url TEXT NOT NULL,
		display_order INT DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS product_sizes (
		id SERIAL PRIMARY KEY,
		product_id VARCHAR(100) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		size VARCHAR(20) NOT NULL,
		stock INT DEFAULT 10,
		UNIQUE(product_id, size)
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(100) PRIMARY KEY,
		customer_name VARCHAR(255) NOT NULL,
		customer_phone VARCHAR(50) NOT NULL,
		customer_email VARCHAR(255) NOT NULL,
		delivery_location TEXT NOT NULL,
		payment_method VARCHAR(50) NOT NULL,
		mpesa_phone VARCHAR(50),
		subtotal NUMERIC(12, 2) NOT NULL,
		delivery_fee NUMERIC(12, 2) NOT NULL,
		total_amount NUMERIC(12, 2) NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id VARCHAR(100) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id VARCHAR(100) NOT NULL,
		size VARCHAR(20) NOT NULL,
		quantity INT NOT NULL,
		price NUMERIC(12, 2) NOT NULL
	);
	`

	_, err := db.Exec(schema)
	return err
}
