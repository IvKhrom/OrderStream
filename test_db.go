package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	dsn := "postgres://postgres:upvel123@localhost:5433/orderstream?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close(context.Background())

	var currentDB string
	err = conn.QueryRow(context.Background(), "SELECT current_database()").Scan(&currentDB)
	if err != nil {
		log.Fatal("Current DB error:", err)
	}
	fmt.Printf("Connected to database: %s\n", currentDB)

	var tableExists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'orders')").Scan(&tableExists)

	if err != nil {
		log.Fatal("Query error:", err)
	}

	fmt.Printf("Table 'orders' exists: %t\n", tableExists)

	if tableExists {
		fmt.Println("SUCCESS! Database is ready!")
	}
}
