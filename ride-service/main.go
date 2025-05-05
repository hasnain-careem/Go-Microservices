package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"ride-service/config"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("❌ Could not connect to DB: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ DB not reachable: %v", err)
	}

	fmt.Println("✅ Connected to rides_db successfully")
}
